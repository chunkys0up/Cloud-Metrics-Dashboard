export type LatencyMetric = {
  date: string;
  ms: number;
};

export type ComputerMetric = {
  resource: "CPU" | "Memory" | "Disk";
  type: "used" | "unused";
  percent: number;
};

export type NetworkMetric = {
  date: string;
  value: number;
  symbol: "RX" | "TX";
};

export type RequestMetric = {
  column_name: string;
  column_value: number;
};

export interface Report {
  Timestamp: string;
  TotalRequests: number;
  FailedRequests: number;
  AverageLatency: number;
  CpuUsed: number;
  MemoryUsed: number;
  DiskUsed: number;
  RxBytesRate: number;
  TxBytesRate: number;
}
