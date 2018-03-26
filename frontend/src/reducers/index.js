import { combineReducers } from 'redux'
import { SWARM_STATE_UPDATE } from '../actions'
import _ from 'lodash'

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

        const orderedTasks = _.orderBy(action.payload.tasks, ['CreatedAt'], ['desc'])

        const tasksByServiceID = _.groupBy(orderedTasks, (t) => t.ServiceID)

        const serviceList = _.map(
          action.payload.services,
          (s) => ({
            name: s.Spec.Name, id: s.ID,
            status: _.map(tasksByServiceID[s.ID], (t)=> t.Status.State)[0],
            createdAt: s.CreatedAt
          })
        )

        return  _.orderBy(serviceList, 'createdAt')
      }
      return state
    }
  }
)

export default rootReducer
