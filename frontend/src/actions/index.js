import axios from 'axios'
const SWARM_STATE_UPDATE='SWARM_STATE_UPDATE'


export function swarmStateUpdate(newState) {
  return {
    type: SWARM_STATE_UPDATE,
    payload: newState
  }
}
