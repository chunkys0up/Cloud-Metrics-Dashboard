import './style.css';
import { LineGraph, BarGraph, StackedBarsGraph, MultiLineChart } from './Metric-Graphs/Graphs';
import { useEffect, useState } from 'react';
import type { Report, ComputerMetric, NetworkMetric } from './redis-data/DataTypes';

const initialRequestMetrics = [
  { column_name: "Total Requests", column_value: 0 },
  { column_name: "Failed Requests", column_value: 0 }
];

const sampleLineData = [
  { date: "2023-01-01T00:00:00", ms: 120 },
  { date: "2023-01-01T00:00:01", ms: 125 },
  { date: "2023-01-01T00:00:02", ms: 118 },
  { date: "2023-01-01T00:00:03", ms: 132 },
  { date: "2023-01-01T00:00:04", ms: 127 },
  { date: "2023-01-01T00:00:05", ms: 140 },
  { date: "2023-01-01T00:00:06", ms: 136 },
  { date: "2023-01-01T00:00:07", ms: 142 },
  { date: "2023-01-01T00:00:08", ms: 138 },
  { date: "2023-01-01T00:00:09", ms: 150 },
  { date: "2023-01-01T00:00:10", ms: 147 },
  { date: "2023-01-01T00:00:11", ms: 145 },
  { date: "2023-01-01T00:00:12", ms: 152 },
  { date: "2023-01-01T00:00:13", ms: 149 },
  { date: "2023-01-01T00:00:14", ms: 155 }
];



function App() {
  const [computerMetrics, setComputerMetrics] = useState<ComputerMetric[]>([]);
  const [requestMetrics, setRequestMetrics] = useState(initialRequestMetrics);
  const [multiLineData, setMultiLineData] = useState<NetworkMetric[]>([]);

  async function getData() {
    try {
      const response = await fetch(`http://localhost:8080/get/`);
      if (!response.ok) throw new Error(`HTTP error, status: ${response.status}`);

      const data: Report = await response.json();

      // console.log("CPU:", data.Metrics.ServerData.CpuUsed);
      // console.log("Memory:", data.Metrics.ServerData.MemoryUsed);
      // console.log("Disk:", data.Metrics.ServerData.DiskUsed);

      // set computer metrics
      const newMetrics: ComputerMetric[] = [
        { resource: "CPU", type: "used", percent: data.Metrics.ServerData.CpuUsed },
        { resource: "CPU", type: "unused", percent: 100 - data.Metrics.ServerData.CpuUsed },
        { resource: "Memory", type: "used", percent: data.Metrics.ServerData.MemoryUsed },
        { resource: "Memory", type: "unused", percent: 100 - data.Metrics.ServerData.MemoryUsed },
        { resource: "Disk", type: "used", percent: data.Metrics.ServerData.DiskUsed },
        { resource: "Disk", type: "unused", percent: 100 - data.Metrics.ServerData.DiskUsed }
      ];
      setComputerMetrics(newMetrics);

      // update the total requests
      setRequestMetrics([
        { column_name: "Total Requests", column_value: data.Metrics.SiteData.TotalRequests },
        { column_name: "Failed Requests", column_value: data.Metrics.SiteData.FailedRequests }
      ]);

      // add to the multilined graph
      const rxMetric: NetworkMetric = {
        date: data.Timestamp,
        value: data.Metrics.ServerData.NetworkTraffic.RxBytesRate,
        symbol: "RX"
      };
      const txMetric: NetworkMetric = {
        date: data.Timestamp,
        value: data.Metrics.ServerData.NetworkTraffic.TxBytesRate,
        symbol: "TX"
      };
      setMultiLineData(prev => [...prev, rxMetric, txMetric]);


    } catch (error) {
      console.error("Error calling API:", error);
    }
  }

  useEffect(() => {
    getData();

    const interval = setInterval(() => {
      getData();
    }, 1000);

    console.log("Went through...");

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="w-screen h-screen flex flex-col ">
      {/* Header */}
      <h1 className="flex-shrink-0 w-full border-b-2 font-bold text-2xl p-3 flex items-center justify-between">
        Cloud Metrics
        <div className="text-sm text-gray-600 italic">Connected to...</div>
      </h1>
  
      {/* Grid Section */}
      <div className="grid grid-rows-2 grid-cols-2 flex-grow gap-4 p-4 min-h-0">
        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <LineGraph data={sampleLineData} />
        </div>
  
        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <BarGraph data={requestMetrics} />
        </div>
  
        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <StackedBarsGraph data={computerMetrics} />
        </div>
  
        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <MultiLineChart data={multiLineData} />
        </div>
      </div>
    </div>
  );
  

}

export default App
