import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'

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
  render() {
    return (
      <div>
          <Link to="/deploy_stack">Deploy Or Update a Stack</Link>
          <div>Services</div>
          <pre>{JSON.stringify(this.props.services,null,2)}</pre>
      </div>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    services: state.services
  }
}


export default connect(mapStateToProps)(Index)
