package extractor

type IndexExtractor interface {
	ListBundlesFromPackage(indexImage string, pkgName string) ([]string, error)
	ListAllBundles(indexImage string) ([]string, error)
}

type IndexExtractorCache interface {
	GetBundles(indexImage string, cacheKey string) []string
	SetBundles(indexImage string, pkgBundlesMap map[string][]string)
}
