package buildops

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"cloud-mta-build-tool/internal/fs"
	"cloud-mta-build-tool/mta"
)

var _ = Describe("BuildParams", func() {

	var _ = Describe("GetBuildResultsPath", func() {
		var _ = DescribeTable("valid cases", func(module *mta.Module, expected string) {
			Ω(GetBuildResultsPath(&dir.Loc{}, module)).Should(HaveSuffix(expected))
		},
			Entry("Implicit Build Results Path", &mta.Module{Path: "mPath"}, "mPath"),
			Entry("Explicit Build Results Path",
				&mta.Module{
					Path:        "mPath",
					BuildParams: map[string]interface{}{buildResultParam: "bPath"},
				},
				"bPath"))
	})

	var _ = DescribeTable("getRequiredTargetPath", func(requires BuildRequires, module mta.Module, expected string) {
		Ω(getRequiredTargetPath(&dir.Loc{}, &module, &requires)).Should(HaveSuffix(expected))
	},
		Entry("Implicit Target Path", BuildRequires{}, mta.Module{Path: "mPath"}, "mPath"),
		Entry("Explicit Target Path", BuildRequires{TargetPath: "artifacts"}, mta.Module{Path: "mPath"}, filepath.Join("mPath", "artifacts")))

	var _ = Describe("ProcessRequirements", func() {

		AfterEach(func() {
			wd, _ := os.Getwd()
			os.RemoveAll(filepath.Join(wd, "testdata", "testproject", "moduleB"))
		})

		var _ = DescribeTable("Valid cases", func(artifacts []string, expectedPath string) {
			wd, _ := os.Getwd()
			ep := dir.Loc{SourcePath: filepath.Join(wd, "testdata", "testproject"), TargetPath: filepath.Join(wd, "testdata", "result")}
			require := BuildRequires{
				Name:       "A",
				Artifacts:  artifacts,
				TargetPath: "./b_copied_artifacts",
			}
			reqs := []BuildRequires{require}
			mtaObj := mta.MTA{
				Modules: []*mta.Module{
					{
						Name: "A",
						Path: "ui5app",
					},
					{
						Name: "B",
						Path: "moduleB",
						BuildParams: map[string]interface{}{
							requiresParam: reqs,
						},
					},
				},
			}
			Ω(ProcessRequirements(&ep, &mtaObj, &require, "B")).Should(Succeed())
			Ω(filepath.Join(wd, expectedPath)).Should(BeADirectory())
			Ω(filepath.Join(wd, expectedPath, "webapp", "Component.js")).Should(BeAnExistingFile())
		},
			Entry("Require All - list", []string{"*"}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All - single value", []string{"*"}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts")),
			Entry("Require All From Parent", []string{"."}, filepath.Join("testdata", "testproject", "moduleB", "b_copied_artifacts", "ui5app")))

		var _ = DescribeTable("Invalid cases", func(lp *dir.Loc, require BuildRequires, mtaObj mta.MTA, moduleName string) {
			Ω(ProcessRequirements(lp, &mtaObj, &require, moduleName)).Should(HaveOccurred())
		},
			Entry("Module not defined",
				&dir.Loc{},
				BuildRequires{Name: "A", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"},
				mta.MTA{Modules: []*mta.Module{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"C"),
			Entry("Required Module not defined",
				&dir.Loc{},
				BuildRequires{Name: "C", Artifacts: []string{"*"}, TargetPath: "b_copied_artifacts"},
				mta.MTA{Modules: []*mta.Module{{Name: "A", Path: "ui5app"}, {Name: "B", Path: "moduleB"}}},
				"B"),
			Entry("Target path - file",
				&dir.Loc{SourcePath: getTestPath("testbuildparams")},
				BuildRequires{Name: "ui1", Artifacts: []string{"*"}, TargetPath: "file.txt"},
				mta.MTA{Modules: []*mta.Module{{Name: "ui1", Path: "ui1"}, {Name: "node", Path: "node"}}},
				"node"))

	})
})

var _ = Describe("Process complex list of requirements", func() {
	AfterEach(func() {
		os.RemoveAll(getTestPath("testbuildparams", "node", "existingfolder", "deepfolder"))
		os.RemoveAll(getTestPath("testbuildparams", "node", "newfolder"))
	})

	It("", func() {
		lp := dir.Loc{
			SourcePath: getTestPath("testbuildparams"),
			TargetPath: getTestPath("result"),
		}
		mtaObj, _ := lp.ParseFile()
		for _, m := range mtaObj.Modules {
			if m.Name == "node" {
				for _, r := range getBuildRequires(m) {
					ProcessRequirements(&lp, mtaObj, &r, "node")
				}
			}
		}
		// ["*"] => "newfolder"
		Ω(getTestPath("testbuildparams", "node", "newfolder", "webapp")).Should(BeADirectory())
		// ["deep/folder/inui2/anotherfile.txt"] => "existingfolder/deepfolder"
		Ω(getTestPath("testbuildparams", "node", "existingfolder", "deepfolder", "anotherfile.txt")).Should(BeAnExistingFile())
		// ["./deep/*/inui2/another*"] => "./existingfolder/deepfolder"
		Ω(getTestPath("testbuildparams", "node", "existingfolder", "deepfolder", "anotherfile2.txt")).Should(BeAnExistingFile())
		// ["deep/folder/inui2/somefile.txt", "*/folder/"] =>  "newfolder/newdeepfolder"
		Ω(getTestPath("testbuildparams", "node", "newfolder", "newdeepfolder", "folder")).Should(BeADirectory())
	})

})

var _ = Describe("PlatformsDefined", func() {
	It("No platforms", func() {
		m := mta.Module{
			Name: "x",
			BuildParams: map[string]interface{}{
				SupportedPlatformsParam: []string{},
			},
		}
		Ω(PlatformsDefined(&m)).Should(Equal(false))
	})
	It("All platforms", func() {
		m := mta.Module{
			Name:        "x",
			BuildParams: map[string]interface{}{},
		}
		Ω(PlatformsDefined(&m)).Should(Equal(true))
	})
})

var _ = Describe("GetBuilder", func() {
	It("Builder defined by type", func() {
		m := mta.Module{
			Name: "x",
			Type: "node-js",
			BuildParams: map[string]interface{}{
				SupportedPlatformsParam: []string{},
			},
		}
		Ω(GetBuilder(&m)).Should(Equal("node-js"))
	})
	It("Builder defined by build params", func() {
		m := mta.Module{
			Name: "x",
			Type: "node-js",
			BuildParams: map[string]interface{}{
				builderParam: "npm",
			},
		}
		Ω(GetBuilder(&m)).Should(Equal("npm"))
	})
})

func getTestPath(relPath ...string) string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "testdata", filepath.Join(relPath...))
}