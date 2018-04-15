import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import filesize from 'filesize'

class Index extends Component {

  constructor(props) {
    super(props)
  }

  state = {
    loadingText: null
  }

  render() {

    const { modal } = this.state

    return (
      <div>
        <div className="container">        
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
