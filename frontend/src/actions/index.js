import axios from 'axios'
export const SWARM_STATE_UPDATE='SWARM_STATE_UPDATE'
export const WEBSOCKET_CONNECTED='WEBSOCKET_CONNECTED'
export const WEBSOCKET_DISCONNECTED='WEBSOCKET_DISCONNECTED'

export function swarmStateUpdate(newState) {
  return {
    type: SWARM_STATE_UPDATE,
    payload: newState
  }
}

export const websocketConnected = () => ({ type: WEBSOCKET_CONNECTED })
export const websocketDisconnected = () => ({ type: WEBSOCKET_DISCONNECTED })
