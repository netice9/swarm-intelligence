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

    const cpuHistory = this.props.swarm.cpuHistory || [new Date(0),0]
    const memoryHistory = this.props.swarm.memoryHistory || [new Date(0),0]

    return (
      <div>
        <div className="container-fluid">
          <div className="row">
            <div className="col">
              <h3>Namespaces</h3>
              <table className="table table-striped table-hover">
                <thead className="thead">
                  <tr>
                    <th>Name</th>
                    <th style={ {textAlign: 'right'} } >Memory Usage</th>
                    <th style={ {textAlign: 'right'} } >% CPU Usage</th>
                  </tr>
                </thead>
                <tbody>
                {
                  _.map(this.props.swarm.namespaces, (ns) => (
                    <tr key={ns.namespace} onClick={ () => {this.props.history.push(`/namespaces/${ns.namespace}`)}}>
                      <td>{ns.namespace}</td>
                      <td style={ {textAlign: 'right'} } >{filesize(ns.memory || 0)}</td>
                      <td style={ {textAlign: 'right'} } >{(ns.cpu * 100).toFixed(2)}</td>
                    </tr>
                  ))
                  }
                </tbody>
              </table>
              </div>
            </div>

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
