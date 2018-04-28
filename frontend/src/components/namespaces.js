import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import filesize from 'filesize'
import { Chart } from 'react-google-charts'

class Index extends Component {

  constructor(props) {
    super(props)
  }

  state = {
    loadingText: null
  }

  render() {

    const { modal } = this.state

    const memoryData = _.map(this.props.swarm.namespaces,(ns) => [ns.namespace,ns.memory])

    const cpuData = _.map(this.props.swarm.namespaces,(ns) => [ns.namespace,ns.cpu*100])

    const cpuHistory = this.props.swarm.cpuHistory || [[new Date(0),0]]
    const memoryHistory = this.props.swarm.memoryHistory || [[new Date(0),0]]

    return (
      <div className="container-fluid">
        <div className="row mt-3">
          <div className="col">
            <div className="card">
              <div className="card-header">
                Namespaces
              </div>
              <div className="card-body">
                <table className="table table-hover table-sm">
                  <thead className="thead thead-light">
                    <tr>
                      <th>Name</th>
                      <th style={ {textAlign: 'right'} } ># Services</th>
                      <th style={ {textAlign: 'right'} } ># Containers</th>
                      <th style={ {textAlign: 'right'} } >Memory Usage</th>
                      <th style={ {textAlign: 'right'} } >% CPU Usage</th>
                    </tr>
                  </thead>
                  <tbody>
                  {
                    _.map(this.props.swarm.namespaces, (ns) => (
                      <tr key={ns.namespace} onClick={ () => {this.props.history.push(`/namespaces/${ns.namespace}`)}}>
                        <td>{ns.namespace}</td>
                        <td style={ {textAlign: 'right'} } >{ns.numberOfServices || 0}</td>
                        <td style={ {textAlign: 'right'} } >{ns.numberOfContainers || 0}</td>
                        <td style={ {textAlign: 'right'} } >{filesize(ns.memory || 0)}</td>
                        <td style={ {textAlign: 'right'} } >{(ns.cpu * 100).toFixed(2)}</td>
                      </tr>
                    ))
                    }
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div className="row mt-3">
          <div className="col">
            <div className="card">
              <div className="card-header">
                Volumes
              </div>
              <div className="card-body">
              <table className="table table-sm">
                <thead className="thead thead-light">
                  <tr>
                    <th>Name</th>
                    <th>Created At</th>
                    <th>Driver</th>
                    <th>Scope</th>

                  </tr>
                </thead>
                <tbody>
                {
                  _.map(this.props.swarm.volumes, (v) => (
                    <tr key={v.Name} >
                      <td>{v.Name}</td>
                      <td>{v.CreatedAt}</td>
                      <td>{v.Driver}</td>
                      <td>{v.Scope}</td>
                    </tr>
                  ))
                  }
                </tbody>
              </table>
              </div>
            </div>
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
                  options={
                    {
                      title: "Memory Usage",
                      legend: {
                        position: 'bottom'
                      }
                    }
                  }

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


        <div className="card  mt-3">
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
                  options={
                    {
                      legend: {
                        position: 'bottom'
                      }
                    }
                  }
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
    )
  }
}

const mapStateToProps = (state) => {
  return {
    swarm: state.swarm,
    swarmState: state.swarmState
  }
}


export default connect(mapStateToProps)(Index)
