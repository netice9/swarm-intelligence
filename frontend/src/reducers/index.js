import { combineReducers } from 'redux'
import { SWARM_STATE_UPDATE, WEBSOCKET_CONNECTED, WEBSOCKET_DISCONNECTED } from '../actions'
import _ from 'lodash'
import moment from 'moment'


const MAX_HISTORY_SIZE=60

const rootReducer = combineReducers(
  {
    websocketConnected: (state = false, action) => {
      switch( action.type ) {
        case WEBSOCKET_CONNECTED:
          return true
        case WEBSOCKET_DISCONNECTED:
          return false
        default:
        return state
      }
    },
    swarmState: (state = {}, action) => {
      if (action.type === 'SWARM_STATE_UPDATE') {
        return {
          ...action.payload
        }
      }
      return state
    },
    swarm: (oldState = {}, action) => {
      if (action.type === 'SWARM_STATE_UPDATE') {

        const orderedTasks = _.orderBy(action.payload.tasks, ['CreatedAt'], ['desc'])

        const tasksByServiceID = _.groupBy(orderedTasks, (t) => t.ServiceID)

        const containersWithStats = _.map(action.payload.containers, (c) => (
            {
              ...c,
              stats: action.payload.stats[c.Id] || {
                memory_stats: {
                  usage: 0
                },
                cpu_stats: {
                  cpu_usage: {
                    total_usage: 0,
                  }
                },
                precpu_stats: {
                  cpu_stats: {
                    cpu_usage: {
                      total_usage: 0,
                    }
                  },
                },
                preread: "2018-03-26T19:56:13.614104943Z",
                read: "2018-03-26T19:56:14.614331085Z"
              }
            }
          )
        )

        const containersByServiceID = _.groupBy(containersWithStats, (c) => (c.Labels['com.docker.swarm.service.id'] || ''))

        const cpuUsage = (serviceID) => {
          const containers = containersByServiceID[serviceID] || []
          const usage = _.sum(_.map(containers, (c) => c.stats.cpu_stats.cpu_usage.total_usage - c.stats.precpu_stats.cpu_usage.total_usage))
          const duration = _.sum(_.map(containers, (c) => moment(c.stats.read).diff(moment(c.stats.preread))))
          if (!usage || !duration) {
            return 0
          }
          return usage / (duration * 1000 * 1000)
        }


        const serviceList = _.map(
          action.payload.services,
          (s) => ({
            name: s.Spec.Name,
            id: s.ID,
            status: _.map(tasksByServiceID[s.ID], (t)=> t.Status.State)[0],
            createdAt: s.CreatedAt,
            memory: _.sum(_.map((containersByServiceID[s.ID] || []),(c) => c.stats.memory_stats.usage)),
            cpu: cpuUsage(s.ID),
            namespace: s.Spec.Labels['com.docker.stack.namespace'] || 'default'
          })
        )

        const byNamespace = _.groupBy(serviceList, 'namespace')

        const namespaces = _.map(byNamespace, (services, ns) => ({
          namespace: ns,
          cpu: _.sumBy(services,'cpu'),
          memory: _.sumBy(services, 'memory'),
          createdAt: _.minBy(services, 'createdAt')
        }))


        const cpu = _.sumBy(namespaces, 'cpu')
        let cpuHistory = [...(oldState.cpuHistory || []), [new Date(Date.parse(action.payload.time)),cpu*100]]

        while (cpuHistory.length > MAX_HISTORY_SIZE) {
          cpuHistory = [...cpuHistory].slice(1)
        }

        const memory = _.sumBy(namespaces, 'memory')
        let memoryHistory = [...(oldState.memoryHistory || []), [new Date(Date.parse(action.payload.time)),memory/(1024*1024)]]

        while (memoryHistory.length > MAX_HISTORY_SIZE) {
          memoryHistory = [...memoryHistory].slice(1)
        }

        const state = {
          cpuHistory,
          memoryHistory,
          cpu: cpu,
          memory: memory,
          namespaces: _.sortBy(namespaces, 'createdAt'),
          servicesByNamespace: byNamespace
        }

        return state
      }
      return oldState
    }
  }
)

export default rootReducer
