import { useEffect, useRef } from "react";
import * as Plot from "@observablehq/plot";

interface LinePoint {
  date: Date;
  ms: number;
}

interface PlotChartProps {
  data: LinePoint[];
}

export function LineGraph({ data }: PlotChartProps) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const chart = Plot.plot({
      width: chartRef.current.clientWidth,
      height: chartRef.current.clientHeight,
      x: { grid: true, label: "Date" },
      y: { grid: true, label: "Average Latency (ms)" },
      marks: [
        Plot.ruleY([0]),
        Plot.lineY(data, { x: "date", y: "ms", stroke: "blue" })
      ]
    });

    chartRef.current.appendChild(chart);

    return () => chart.remove();
  }, [data]);

  return (
    <div
      ref={chartRef}
      className="w-full h-full"
      style={{ minHeight: "300px" }}
    ></div>
  );
}


interface BarPoint {
  column_name: string;
  column_value: number;
}

interface BarGraphProps {
  data: BarPoint[];
}

export function BarGraph({ data }: BarGraphProps) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const chart = Plot.plot({
      width: chartRef.current.clientWidth,
      height: chartRef.current.clientHeight,
      marks: [
        Plot.axisX({ label: null, lineWidth: 8 }),
        Plot.axisY({ label: "Requests" }),
        Plot.ruleY([0]),
        Plot.barY(data, { x: "column_name", y: "column_value", fill: "steelblue" })
      ]
    });

    chartRef.current.appendChild(chart);
    return () => chart.remove();
  }, [data]);

  return (
    <div
      ref={chartRef}
      className="w-full h-full"
      style={{ minHeight: "300px" }}
    ></div>
  );
}


interface StackedBarPoint {
  resource: string;
  type: string;
  percent: number;
}

interface StackedBarProps {
  data: StackedBarPoint[];
}

export function StackedBarsGraph({ data }: StackedBarProps) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const chart = Plot.plot({
      marginLeft: 90,
      x: { label: "Usage (%)", grid: true },
      color: { legend: true },

      marks: [
        Plot.barX(
          data,
          Plot.groupY(
            { x: "sum" },
            { y: "resource", x: "percent", fill: "type", z: "type" }
          )
        )
      ]

    });

    chartRef.current.appendChild(chart);
    return () => chart.remove();
  }, [data]);

  return (
    <div
      ref={chartRef}
      className="w-full h-full"
      style={{ maxHeight: "300px" }}
    ></div>
  );
}


interface MultiLinePoint {
  date: Date;
  value: number;
  symbol: string;
}

interface MultiLineProps {
  data: MultiLinePoint[];
}

export function MultiLineChart({ data }: MultiLineProps) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const chart = Plot.plot({
      style: "overflow: visible;",
      y: { grid: true },
      width: chartRef.current.clientWidth,
      height: chartRef.current.clientHeight,
      marks: [
        Plot.ruleY([0]),
        Plot.lineY(data, { x: "date", y: "value", stroke: "symbol" }),
        Plot.text(data, Plot.selectLast({ x: "date", y: "value", z: "symbol", text: "symbol", textAnchor: "start", dx: 3 }))
      ]
    });

    chartRef.current.appendChild(chart);
    return () => chart.remove();
  }, [data]);

  return (
    <div
      ref={chartRef}
      className="w-full h-full"
      style={{ maxHeight: "300px" }}
    ></div>
  );
}




