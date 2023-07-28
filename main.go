package main

import (
	_ "embed"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed framego.go.tmpl
var framegoTpl string

type FramegoModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

func NewFramegoModule() *FramegoModule {
	return &FramegoModule{ModuleBase: &pgs.ModuleBase{}}
}

func (p *FramegoModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	tpl := template.New("framego").Funcs(map[string]interface{}{
		"Package": p.ctx.PackageName,
		"Name":    p.ctx.Name,
	})

	p.tpl = template.Must(tpl.Parse(framegoTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *FramegoModule) Name() string {
	return "framegoTpl"
}

func (p *FramegoModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		p.generate(t)
	}

	return p.Artifacts()
}

func (p *FramegoModule) generate(f pgs.File) {
	if len(f.Services()) == 0 {
		return
	}
	names := strings.Split(f.Services()[0].FullyQualifiedName(), ".")
	packageName := names[len(names)-2]
	services := make([]string, 0)
	for _, service := range f.Services() {
		services = append(services, string(service.Name()))
	}
	data := map[string]interface{}{
		"PackageName": packageName,
		"Services":    services,
	}

	name := p.ctx.OutputPath(f).SetExt(".framego.go")
	p.AddGeneratorTemplateFile(name.String(), p.tpl, data)
}

func main() {
	feature := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.
		Init(pgs.DebugEnv("DEBUG"), pgs.SupportedFeatures(&feature)).
		RegisterModule(NewFramegoModule()).
		RegisterPostProcessor(pgsgo.GoFmt()).
		Render()
}
