import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import filesize from 'filesize'
import removeIcon from '../icons/ic_remove_circle_red_18px.svg'
import viewIcon from '../icons/ic_visibility_black_18px.svg'
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap'
import axios from 'axios';
import Loadable from 'react-loading-overlay'
import { Chart } from 'react-google-charts'
import LogTailer from './log_tailer'

class Services extends Component {

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

  toggleModal() {
    this.setState({modal: null})
  }

  render() {

    const { modal } = this.state
    const namespaceName = this.props.match.params.namespace
    const namespace = _.find((this.props.swarm.namespaces || []), ns => namespaceName === ns.namespace)

    if (!namespace) {
      return "Loading ..."
    }

    const services = namespace.services

    const memoryData = _.map(services,(s) => [s.name,s.memory])
    const cpuData = _.map(services,(s) => [s.name,s.cpu*100])

    const {cpuHistory, memoryHistory} = namespace


    return (
      <Loadable
        spinner
        active={!!this.state.loadingText}
        text={this.state.loadingText}
      >
        {
          this.state.modal &&
          <Modal isOpen={true} toggle={this.toggleModal} size="lg">
            <ModalHeader toggle={this.toggleModal}>{modal.title}</ModalHeader>
            <ModalBody>
              { modal.text }
            </ModalBody>
            <ModalFooter>
              { modal.confirmText && <Button color="warning" onClick={modal.confirmAction}>{modal.confirmText}</Button> }
              {' '}
              { modal.showCancel && <Button color="secondary" onClick={this.toggleModal}>Cancel</Button> }
            </ModalFooter>
          </Modal>
        }

        <div className="container-fluid">
          <div className="card mt-3">
            <div className="card-header">
              Services for {namespaceName}
            </div>
            <div className="card-body">
              <table className="table table-hover table-sm">
                <thead className="thead thead-light">
                  <tr>
                    <th>Name</th>
                    <th>Status</th>
                    <th style={ {textAlign: 'right'} } >Memory Usage</th>
                    <th style={ {textAlign: 'right'} } >% CPU Usage</th>
                    <th>Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {
                    _.map(
                       services,
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
                          <img src={viewIcon} onClick={() =>{
                            this.setState({
                              modal: {
                                title: `Logs for ${s.name}`,
                                text: (<LogTailer url={`/api/services/${s.id}/logs`} />),
                                confirmText: 'Close'
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
          </div>

          <div className="card mt-3">
            <div className="card-header">
              Memory
            </div>
            <div className="card-body">
              <div className="row">
                <div className="col-3">
                  <Chart
                    chartType="PieChart"
                    data={[
                      ['Namespace', 'Bytes Used'],
                      ...memoryData
                    ]}
                    options={{
                      legend: {
                        position: 'bottom'
                      }
                    }}

                    width="100%"
                    graph_id="MemoryChart"
                  />
                </div>

                <div className="col-9">
                  <Chart
                    chartType="LineChart"
                    columns={[
                      {
                        type: 'datetime'
                      },
                      {
                        label: 'Mbytes',
                        type: 'number'
                      },

                    ]}
                    rows={memoryHistory}
                    options={
                      {
                        title: "Memory History",
                        hAxis: {
                           format: "HH:mm:ss",
                           title: 'Time',
                           slantedText: true
                        },
                        vAxis: {
                          baseline: 0
                        }
                      }
                    }
                    width="100%"
                    graph_id="MemoryHistoryChart"
                  />
                </div>

              </div>
            </div>
          </div>

          <div className="card mt-3">
            <div className="card-header">
              CPU
            </div>
            <div className="card-body">
              <div className="row">
                <div className="col-3">
                  <Chart
                    chartType="PieChart"
                    data={[
                      ['Namespace', '%CPU'],
                      ...cpuData
                    ]}
                    options={{
                      legend: {
                        position: 'bottom'
                      }
                    }}
                    width="100%"
                    graph_id="CPUChart"
                  />
                </div>
                <div className="col-9">
                  <Chart
                    chartType="LineChart"
                    columns={[
                      {
                        type: 'datetime'
                      },
                      {
                        label: '%',
                        type: 'number'
                      },

                    ]}
                    rows={cpuHistory}
                    options={
                      {
                        hAxis: {
                           format: "HH:mm:ss",
                           title: 'Time',
                           slantedText: true
                        },
                        vAxis: {
                          baseline: 0
                        }
                      }
                    }
                    width="100%"
                    graph_id="CPUHistoryChart"
                  />
                </div>
              </div>
            </div>
          </div>

        </div>

      </Loadable>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    swarm: state.swarm
  }
}


export default connect(mapStateToProps)(Services)
