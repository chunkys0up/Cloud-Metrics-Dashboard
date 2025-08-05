import './style.css';
import { LineGraph, BarGraph, StackedBarsGraph, MultiLineChart } from './Metric-Graphs/Traffic';

function App() {

  const sampleLineData = [
    { date: new Date("2023-01-01"), ms: 120 },
    { date: new Date("2023-01-02"), ms: 130 },
    { date: new Date("2023-01-03"), ms: 125 },
    { date: new Date("2023-01-04"), ms: 135 },
    { date: new Date("2023-01-05"), ms: 128 },
    { date: new Date("2023-01-06"), ms: 140 },
    { date: new Date("2023-01-07"), ms: 145 },
    { date: new Date("2023-01-08"), ms: 138 },
    { date: new Date("2023-01-09"), ms: 150 },
    { date: new Date("2023-01-10"), ms: 155 },
    { date: new Date("2023-01-11"), ms: 90 },
    { date: new Date("2023-01-12"), ms: 75 },
    { date: new Date("2023-01-13"), ms: 83 },
    { date: new Date("2023-01-14"), ms: 101 },
  ];

  const sampleBarData = [
    { column_name: "Total Requests", column_value: 123 },
    { column_name: "Failed Requests", column_value: 4 }
  ];

  const sampleStackedBarsData = [
    { resource: "CPU", type: "used", percent: 23 },
    { resource: "CPU", type: "unused", percent: 77 },
    { resource: "Memory", type: "used", percent: 45 },
    { resource: "Memory", type: "unused", percent: 55 },
    { resource: "Disk", type: "used", percent: 5 },
    { resource: "Disk", type: "unused", percent: 95 }
  ];

  const multiDataLineData = [
    { date: new Date("2025-08-04T12:00:00"), value: 1200, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:00"), value: 800, symbol: "TX" },

    { date: new Date("2025-08-04T12:00:05"), value: 1500, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:05"), value: 950, symbol: "TX" },

    { date: new Date("2025-08-04T12:00:10"), value: 1700, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:10"), value: 1100, symbol: "TX" },

    { date: new Date("2025-08-04T12:00:15"), value: 1400, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:15"), value: 1050, symbol: "TX" },

    { date: new Date("2025-08-04T12:00:20"), value: 1600, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:20"), value: 1200, symbol: "TX" },

    { date: new Date("2025-08-04T12:00:25"), value: 1900, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:25"), value: 1400, symbol: "TX" },

    { date: new Date("2025-08-04T12:00:30"), value: 2100, symbol: "RX" },
    { date: new Date("2025-08-04T12:00:30"), value: 1600, symbol: "TX" }
  ];


  return (
    <div className="w-screen h-screen flex flex-col ">
      {/* Header */}
      <h1 className="w-full border-b-2 font-bold text-2xl p-3 flex items-center justify-between">
        Cloud Metrics

        <div className="text-sm text-gray-600 italic">Connected to...</div>
      </h1>

      {/* Grid Section */}
      <div className="grid grid-rows-2 grid-cols-2 flex-grow p-1 gap-1">

        <div className="bg-white w-[95%] h-[95%] border-2 border-gray-300 p-2 rounded-3xl">
          <LineGraph data={sampleLineData} />
        </div>

        <div className="bg-white w-[95%] h-[95%] border-2 border-gray-300 p-2 rounded-3xl">
          <BarGraph data={sampleBarData} />
        </div>

        <div className="flex justify-center items-center bg-white w-[95%] h-[95%] border-2 border-gray-300 p-2 rounded-3xl">
          <StackedBarsGraph data={sampleStackedBarsData} />
        </div>

        <div className="flex justify-center items-center bg-white w-[95%] h-[95%] border-2 border-gray-300 p-2 rounded-3xl">
          <MultiLineChart data={multiDataLineData} />
        </div>

      </div>
    </div>

  )
}

export default App
