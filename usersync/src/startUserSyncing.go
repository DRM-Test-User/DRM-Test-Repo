package usersync

import (
	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

type UserSync struct {
	CommitHash  string
	AuthorEmail string
	RepoUrl     string
}

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	ghAccessToken string,
	batchSize int,
	githubGraphQLUrl string,
) {
	newCommitAuthorsRaw, err := getNewCommitAuthors(db)
	if err != nil {
		logger.LogFatalRedAndExit("error getting new commit authors to process: %s", err)
		return
	}
	if newCommitAuthorsRaw == nil {
		logger.LogBlue("no new authors to sync")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := convertDatabaseObjectToUserSync(newCommitAuthorsRaw)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := generateBatchAuthors(repoUrlToAuthorsMap, 2)

	// Get info for each batch

	for _, repoToAuthorBatch := range repoToAuthorBatches {
		logger.LogGreenDebug("%s", repoToAuthorBatch.RepoURL)

		githubGraphQLCommitAuthorsMap, err := identifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.AuthorCommitTuples, ghAccessToken, githubGraphQLUrl)
		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		if githubGraphQLCommitAuthorsMap == nil {
			logger.LogError("commits is nil")
			continue
		}

		githubGraphQLCommitAuthors := make([]GithubGraphQLCommit, 0, len(githubGraphQLCommitAuthorsMap))

		for _, commitAuthor := range githubGraphQLCommitAuthorsMap {
			githubGraphQLCommitAuthors = append(githubGraphQLCommitAuthors, commitAuthor)
		}

		for _, commitAuthor := range githubGraphQLCommitAuthors {
			author := commitAuthor.Author

			err := insertIntoRestIdToUser(author, db)
			if err != nil {
				logger.LogError("error occured while inserting RestID to Email: %s", err)
			}

			err = insertGithubUser(author, db)
			if err != nil {
				logger.LogError("error occured while inserting author: %s", err)
			}

		}
	}
}
