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
        const serviceByID = {}
        action.payload.services.forEach((s) => {
          serviceByID[s.ID]={name: s.Spec.Name, id: s.ID, tasks: []}
        })

        action.payload.tasks.forEach((t) => {
          const s = serviceByID[t.ServiceID] || { tasks: [] }
          s.tasks.push({id: t.ID, createdAt: t.CreatedAt, state: t.Status.State})
        })

        return action.payload.services.map((s) => (serviceByID[s.ID]) )
      }
      return state
    }
  }
)

export default rootReducer
