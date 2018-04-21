import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import filesize from 'filesize'
import removeIcon from '../icons/ic_remove_circle_red_18px.svg'
import { Button, Modal, ModalHeader, ModalBody, ModalFooter } from 'reactstrap'
import axios from 'axios';
import Loadable from 'react-loading-overlay'
import { Chart } from 'react-google-charts'

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
        {
          <div className="container-fluid">
            <h3>Namespace: {namespaceName}</h3>
            <table className="table table-striped table-hover">
              <thead className="thead">
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
                        </td>
                      </tr>
                    )
                  )
                }
              </tbody>
            </table>

            <div className="row">
              <div className="col">
                <Chart
                  chartType="PieChart"
                  data={[
                    ['Namespace', 'Bytes Used'],
                    ...memoryData
                  ]}
                  options={ {title: "Memory Usage"} }
                  pieHole={0.4}
                  width="100%"
                  graph_id="MemoryChart"
                />
              </div>

              <div className="col">
                <Chart
                  chartType="PieChart"
                  data={[
                    ['Namespace', '%CPU'],
                    ...cpuData
                  ]}
                  options={ {title: "CPU Usage"} }
                  width="100%"
                  pieHole={0.4}
                  graph_id="CPUChart"
                />
              </div>
            </div>

            <div className="row">
              <div className="col">
              <Chart
                chartType="LineChart"
                columns={[
                  {
                    label: 'time',
                    type: 'datetime'
                  },
                  {
                    label: 'memory (Mbytes)',
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
              <div className="col">
              <Chart
                chartType="LineChart"
                columns={[
                  {
                    label: 'time',
                    type: 'datetime'
                  },
                  {
                    label: 'cpu (%)',
                    type: 'number'
                  },

                ]}
                rows={cpuHistory}
                options={
                  {
                    title: "CPU History",
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
        }

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
