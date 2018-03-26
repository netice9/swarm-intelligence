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

        const containersWithStats = _.map(action.payload.containers, (c) => (
            {
              ...c,
              stats: action.payload.stats[c.Id] || {
                memory_stats: {
                  usage: 0
                }
              }
            }
          )
        )

        const containersByServiceID = _.groupBy(containersWithStats, (c) => (c.Labels['com.docker.swarm.service.id'] || ''))

        const serviceList = _.map(
          action.payload.services,
          (s) => ({
            name: s.Spec.Name,
            id: s.ID,
            status: _.map(tasksByServiceID[s.ID], (t)=> t.Status.State)[0],
            createdAt: s.CreatedAt,
            memory: _.sum(_.map((containersByServiceID[s.ID] || []),(c) => c.stats.memory_stats.usage))
          })
        )

        return  _.orderBy(serviceList, 'createdAt')
      }
      return state
    }
  }
)

export default rootReducer
