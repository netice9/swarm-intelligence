import React, { Component } from 'react'
import { Link } from 'react-router-dom'

class NavBar extends Component {
  render() {
    return (
      <nav className="navbar navbar-expand-lg navbar-light bg-light sticky-top">
        <Link to="/" className="navbar-brand">Swarm Intelligence</Link>

        <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span className="navbar-toggler-icon"></span>
        </button>

        <div className="collapse navbar-collapse" id="navbarSupportedContent">
          <Link className="nav-link" to="/deploy_stack">Deploy Or Update a Stack</Link>
          <Link className="nav-link" to="/manage_credentials">Manage Credentials</Link>
      </div>

      </nav>
    )
  }
}

export default NavBar
