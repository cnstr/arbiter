package database

// TODO:
// SELECT "architecture", "package", "version", "repositoryTier", "repositorySlug"
// FROM "Package"
// WHERE "isPruned" = false
// AND "isCurrent" = true
// AND "package" IN (
// 	SELECT "package"
// 	FROM "Package"
// 	WHERE "isPruned" = false
// 	AND "isCurrent" = true
// 	GROUP BY "package"
// 	HAVING COUNT(*) > 1
// )
// AND "architecture" NOT IN (
// 	SELECT "architecture"
// 	FROM "Package"
// 	WHERE "isPruned" = false
// 	AND "isCurrent" = true
// 	GROUP BY "architecture"
// 	HAVING COUNT(*) > 1
// );
//
