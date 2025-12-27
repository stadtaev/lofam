import { getRouter } from './router'
import { getRouterManifest } from '@tanstack/react-start/router-manifest'
import { createStartHandler, defaultStreamHandler } from '@tanstack/react-start/server'

export default createStartHandler({
  getRouter,
  getRouterManifest,
})(defaultStreamHandler)
