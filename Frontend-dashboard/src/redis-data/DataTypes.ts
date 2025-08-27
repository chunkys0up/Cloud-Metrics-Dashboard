export type ComputerMetric = {
  resource: "CPU" | "Memory" | "Disk";
  type: "used" | "unused";
  percent: number;
};

export type NetworkMetric = {
  date: string;
  value: number;
  symbol: "RX" | "TX";
}

export interface SiteData { 
  TotalRequests: number;
  FailedRequests: number;
  AverageLatencyMs: number[];
}
export interface NetworkTraffic {
  RxBytesRate: number;
  TxBytesRate: number;
}
export interface ServerData {
  CpuUsed: number;
  MemoryUsed: number;
  DiskUsed: number;
  NetworkTraffic: NetworkTraffic;
}
export interface Metrics {
  SiteData: SiteData;
  ServerData: ServerData;
}
export interface Report {
  Timestamp: string;
  Metrics: Metrics;
}