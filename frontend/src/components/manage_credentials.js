import React, { Component } from 'react'
import { post } from 'axios';
import Loadable from 'react-loading-overlay'

export default class ManageCredentials extends Component {

  state = {
    registry: "",
    username: "",
    password: "",
    alert: null,
    loadingText: null
  }

  canSubmit() {
    return this.state.registry !== "" && this.state.username !== "" && this.state.password !== ""
  }

  changeRegistry = (evt) => {
    this.setState({registry: evt.target.value})
  }

  changeUsername = (evt) => {
    this.setState({username: evt.target.value})
  }

  changePassword = (evt) => {
    this.setState({password: evt.target.value})
  }

  submit = (evt) => {

    this.setState({loadingText: `Adding Credentials`})

    evt.preventDefault()
    post('/api/add_credentials', {
        registry: this.state.registry,
        username: this.state.username,
        password: this.state.password
      }
    ).then(() => {
      this.setState({loadingText: null})
      this.props.history.push("/")
    }).catch((err) => {
      this.setState({loadingText: null})
      if (err.response) {
        this.setState({alert: err.response.data})
      } else {
        this.setState({alert: 'Failed to add credentials'})
      }
    })
  }

  render() {
    return (
      <Loadable
        spinner
        text={this.state.loadingText}
        active={!!this.state.loadingText}
      >
        {this.state.alert ? <div class="alert alert-danger" role="alert">{this.state.alert}</div> : null}
        <div className="container">
          <form onSubmit={this.submit}>
            <div className="form-group">
              <label htmlFor="registry">Registry address</label>
              <input type="text" className="form-control" id="registry" aria-describedby="registryDesc" placeholder="Enter registry address" value={this.state.registry} onChange={this.changeRegistry} />
              <small id="registryDesc" className="form-text text-muted">Name of the stack to be deployed or updated</small>
            </div>
            <div className="form-group">
              <label htmlFor="username">Username</label>
              <input type="text" className="form-control" id="username" aria-describedby="usernameDesc" placeholder="Enter username" value={this.state.username} onChange={this.changeUsername} />
              <small id="usernameDesc" className="form-text text-muted">Username for the registry</small>
            </div>
            <div className="form-group">
              <label htmlFor="password">Password</label>
              <input type="password" className="form-control" id="password" aria-describedby="passwordDesc" placeholder="Enter password" value={this.state.password} onChange={this.changePassword} />
              <small id="passwordDesc" className="form-text text-muted">Password for the registry</small>
            </div>
            <button type="submit" className="btn btn-primary" disabled={!this.canSubmit()}>Submit</button>
          </form>
        </div>
      </Loadable>
    )
  }
}
