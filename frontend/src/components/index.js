import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import _ from 'lodash'
import filesize from 'filesize'

class Index extends Component {

  componentDidMount() {
    this.props.fetchStats
    this.intervalID = setInterval(this.props.fetchStats, 1000)
  }
  componentWillUnmount() {
    if (this.intervalID) {
      clearInterval(this.intervalID)
    }
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

  render() {
    return (
      <div>
        <Link to="/deploy_stack">Deploy Or Update a Stack</Link>
        <div className="container">
            <h3>Services</h3>
            <ul className="list-group list-group-flush">
              {
                _.map(
                   this.props.services,
                  (s) =>(
                    <li key={s.id} className={`d-flex justify-content-between align-items-center ${this.serviceListItemClass(s)}`}>
                      {s.name}
                      <span className="badge badge-info badge-pill">{s.status}</span>
                      <p>{filesize(s.memory)}</p>
                    </li>
                  )
                )
              }
            </ul>
        </div>
      </div>
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
