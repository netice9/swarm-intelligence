import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import filesize from 'filesize'
import removeIcon from '../icons/ic_remove_circle_red_18px.svg'
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap'
import axios from 'axios';
import Loadable from 'react-loading-overlay'

class Index extends Component {

  constructor(props) {
    super(props)
    this.toggleModal = this.toggleModal.bind(this)
    this.deleteService = this.deleteService.bind(this)
  }

  state = {
    selectedService: null,
    modal: null,
    loadingText: null
  }

  deleteService(id, name) {
    this.setState({loadingText: 'Deleting', modal: null})
    axios.delete(`/api/services/${id}`)
      .then((response) => {
        this.setState({
          loadingText: null,
          modal: {
            title: 'Success',
            text: `Service ${name} deleted!`,
            confirmText: 'Ok',
            showCancel: false,
            confirmAction: this.toggleModal
          }
        })
      })
      .catch( (error) => {
        this.setState({
          loadingText: null,
          modal: {
            title: 'Failed',
            text: `Could not delete service ${name}: ${error.response.data}`,
            confirmText: 'Ok',
            showCancel: false,
            confirmAction: this.toggleModal
          }
        })
      })
  }

  serviceListItemClass(s) {
    switch(s.status) {
      case 'running':
        return 'list-group-item list-group-item-success'
      case 'ready':
        return 'list-group-item list-group-item-warning'
      case 'starting':
        return 'list-group-item list-group-item-info'
      case 'complete':
        return 'list-group-item list-group-item-primary'
      case 'shutdown':
        return 'list-group-item list-group-item-dark'
      default:
        return 'list-group-item'
    }
  }

  toggleModal() {
    this.setState({modal: null})
  }

  render() {

    const { modal } = this.state

    return (
      <Loadable
        spinner
        active={!!this.state.loadingText}
        text={this.state.loadingText}
      >
        <Link to="/deploy_stack">Deploy Or Update a Stack</Link>
        <div>
          {
            this.state.modal ?
              <Modal isOpen={true} toggle={this.toggleModal}>
                <ModalHeader toggle={this.toggleModal}>{modal.title}</ModalHeader>
                <ModalBody>
                  { modal.text }
                </ModalBody>
                <ModalFooter>
                  <Button color="warning" onClick={modal.confirmAction}>{modal.confirmText}</Button>{' '}
                  { modal.showCancel ? <Button color="secondary" onClick={this.toggleModal}>Cancel</Button> : null }
                </ModalFooter>
              </Modal>
              :
              null
        }
        </div>
        <div className="container">
            <h3>Services</h3>
            <table className="table table-striped table-hover">
              <thead className="thead">
                <tr>
                  <th>Name</th>
                  <th>Status</th>
                  <th style={ {textAlign: 'right'} } >Memory Usage</th>
                  <th style={ {textAlign: 'right'} } >CPU Usage</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {
                  _.map(
                     this.props.services,
                    (s) =>(
                      <tr key={s.id}>
                        <td>{s.name}</td>
                        <td><span className="badge badge-info badge-pill">{s.status}</span></td>
                        <td style={ {textAlign: 'right'} } >{filesize(s.memory || 0)}</td>
                        <td style={ {textAlign: 'right'} } >{(s.cpu * 100).toFixed(2)}</td>
                        <td>
                          <img src={removeIcon} onClick={() =>{
                            this.setState({
                              modal: {
                                title: 'Delete Service Confirmation',
                                text: `Do you really want to delete service ${s.name}?`,
                                confirmText: 'Delete!',
                                showCancel: true,
                                confirmAction: () => this.deleteService(s.id, s.name)
                              }
                            }
                        )} } />
                        </td>
                      </tr>
                    )
                  )
                }
              </tbody>
          </table>
        </div>
      </Loadable>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    services: state.services,
    swarmState: state.swarmState
  }
}


export default connect(mapStateToProps)(Index)
