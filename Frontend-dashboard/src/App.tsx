import './style.css';
import { LineGraph, BarGraph, StackedBarsGraph, MultiLineChart } from './Metric-Graphs/Graphs';
import { useEffect, useState } from 'react';
import type { Report, LatencyMetric, ComputerMetric, NetworkMetric, RequestMetric } from './redis-data/DataTypes';

function App() {
  const [latencyMetrics, setLatencyMetrics] = useState<LatencyMetric[]>([]);
  const [computerMetrics, setComputerMetrics] = useState<ComputerMetric[]>([]);
  const [requestMetrics, setRequestMetrics] = useState<RequestMetric[]>([]);
  const [byteMetrics, setByteMetrics] = useState<NetworkMetric[]>([]);

  async function getData() {
    try {
      const response = await fetch(`http://localhost:8080/getData`);
      if (!response.ok) throw new Error(`HTTP error, status: ${response.status}`);

      const data: Report = await response.json();

      // add to latency metrics
      const LatencyMetric: LatencyMetric = {
        date: data.Timestamp,
        ms: data.AverageLatency
      }

      setLatencyMetrics(prev => {
        const newArray = [...prev, LatencyMetric];

        if (newArray.length > 200)
          return newArray.slice(-200);

        return newArray;
      });

      // set computer metrics
      const newMetrics: ComputerMetric[] = [
        { resource: "CPU", type: "used", percent: data.CpuUsed },
        { resource: "CPU", type: "unused", percent: 100 - data.CpuUsed },
        { resource: "Memory", type: "used", percent: data.MemoryUsed },
        { resource: "Memory", type: "unused", percent: 100 - data.MemoryUsed },
        { resource: "Disk", type: "used", percent: data.DiskUsed },
        { resource: "Disk", type: "unused", percent: 100 - data.DiskUsed }
      ];

      setComputerMetrics(newMetrics);

      // update the total requests
      const newRequestMetrics: RequestMetric[] = [
        { column_name: "Total Requests", column_value: data.TotalRequests },
        { column_name: "Failed Requests", column_value: data.FailedRequests }
      ];
      setRequestMetrics(newRequestMetrics);

      // add to the multilined graph
      const rxMetric: NetworkMetric = {
        date: data.Timestamp,
        value: data.RxBytesRate,
        symbol: "RX"
      };
      const txMetric: NetworkMetric = {
        date: data.Timestamp,
        value: data.TxBytesRate,
        symbol: "TX"
      };

      setByteMetrics(prev => {
        const tempMetrics = [...prev, rxMetric, txMetric];

        if (tempMetrics.length > 200) {
          return tempMetrics.slice(-200)
        }

        return tempMetrics;
      });

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
      </h1>

      {/* Grid Section */}
      <div className="grid grid-rows-2 grid-cols-2 flex-grow gap-4 p-4 min-h-0">
        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <LineGraph data={latencyMetrics} />
        </div>

        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <BarGraph data={requestMetrics} />
        </div>

        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <StackedBarsGraph data={computerMetrics} />
        </div>

        <div className="bg-white border-2 border-gray-300 p-4 rounded-3xl flex flex-col min-h-0">
          <MultiLineChart data={byteMetrics} />
        </div>
      </div>
    </div>
  );


}

export default App
