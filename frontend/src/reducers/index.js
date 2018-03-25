import { combineReducers } from 'redux'
import { SWARM_STATE_UPDATE } from '../actions'
const rootReducer = combineReducers(
  {
    swarmState: (state = {}, action) => {
      if (action.type === 'SWARM_STATE_UPDATE') {
        return {
          ...action.payload
        }
      }
      return state
    },
    services: (state = [], action) => {
      if (action.type === 'SWARM_STATE_UPDATE') {
        return action.payload.services.map((service) => service.Spec.Name)
      }
      return state
    }
  }
)

export default rootReducer
