import type { Plugin } from "@opencode-ai/plugin"
// Legacy discovery name intentionally inert; bootstrap owns registration.
const adapter: Plugin = async () => ({})
export default adapter
