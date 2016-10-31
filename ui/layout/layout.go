package layout

import "gitlab.netice9.com/dragan/go-reactor/core"

func WithLayout(view *core.DisplayModel) *core.DisplayModel {
	withLayout := layoutUI.DeepCopy()
	withLayout.ReplaceChild("content", view)
	return withLayout
}

var layoutUI = core.MustParseDisplayModel(`
  <div>
  	<bs.Navbar bool:fluid="true">
  		<bs.Navbar.Header>
  			<bs.Navbar.Brand>
  				<a href="#/" className="navbar-brand">Swarm Intelligence</a>
  			</bs.Navbar.Brand>
  		</bs.Navbar.Header>
  	</bs.Navbar>


  	<bs.Grid bool:fluid="true">
			<bs.Row>
				<bs.Col int:mdOffset="1" int:md="10" int:smOffset="0" int:sm="12">
					<div id="content" className="container"/>
				</bs.Col>
			</bs.Row>
  	</bs.Grid>
  </div>

`)
