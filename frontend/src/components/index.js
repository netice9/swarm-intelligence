import React, { Component } from 'react'
import { connect } from 'react-redux'
import { fetchStats } from '../actions'
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
      </div>
    )
  }
}

const mapStateToProps = (state) => {
  return {
    stats: state.stats
  }
}


export default connect(mapStateToProps, {fetchStats} )(Index)
