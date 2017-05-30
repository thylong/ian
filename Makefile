doc-publish:
	hugo -t hugo-theme-learn -s docs-source -d ../docs

doc-serve:
	hugo server --buildDrafts -t hugo-theme-learn -s docs-source -w
