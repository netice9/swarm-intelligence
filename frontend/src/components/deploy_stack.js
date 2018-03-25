import React, { Component } from 'react'
import { post } from 'axios';

export default class DeployStack extends Component {

  state = {
    file: null,
    name: "",
    alert: null,
  }

  canSubmit() {
    return this.state.name !== "" && this.state.file !== null
  }

  changeName = (evt) => {
    this.setState({name: evt.target.value})
  }

  changeFile = (evt) => {
    this.setState({file: evt.target.files[0], alert: null})
  }

  submit = (evt) => {
    evt.preventDefault()
    this.fileUpload().then(() => {
      this.props.history.push("/")
    }).catch((err) => {
      // console.log("catch!")
      if (err.response) {
        this.setState({alert: err.response.data})
      }
    })
  }

  fileUpload() {
    const url = '/api/deploy_stack'
    const formData = new FormData()
    formData.append('file',this.state.file)
    formData.append('name',this.state.name)
    const config = {
      headers: { 'content-type': 'multipart/form-data' }
    }
    return post(url, formData,config)
  }

  render() {
    return (
      <div>
      {this.state.alert ? <div class="alert alert-danger" role="alert">{this.state.alert}</div> : null}

        <div className="container">
          <form onSubmit={this.submit}>
          <div className="form-group">
             <label htmlFor="stack_name">Stack name</label>
             <input type="text" className="form-control" id="stack_name" aria-describedby="stackName" placeholder="Enter Stack Name" value={this.state.name || ""} onChange={this.changeName} />
             <small id="stackName" className="form-text text-muted">Name of the stack to be deployed or updated</small>
           </div>
           <div className="form-group">
             <label htmlFor="composeFile">Compose File</label>
             <input type="file" className="form-control" id="composeFile" placeholder="Compose File" onChange={this.changeFile} />
           </div>
           <button type="submit" className="btn btn-primary" disabled={!this.canSubmit()}>Submit</button>
         </form>
        </div>
      </div>
    )
  }
}
