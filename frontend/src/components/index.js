import React, { Component } from 'react'
import { connect } from 'react-redux'
import { fetchStats } from '../actions'
import moment from 'moment'

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
        Swarm Intelligence
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
