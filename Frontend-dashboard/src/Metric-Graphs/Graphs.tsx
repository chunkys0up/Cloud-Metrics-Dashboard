import { useEffect, useRef } from "react";
import * as Plot from "@observablehq/plot";

export function LineGraph({ data }: any) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const dataMax = Math.max(...data.map((d: any) => d.ms));
    const defaultMax = 10;

    const xMax = Math.max(defaultMax, dataMax);

    const chart = Plot.plot({
      width: chartRef.current.clientWidth,
      height: chartRef.current.clientHeight,
      x: { grid: true, label: "Date" },
      y: { grid: true, label: "Average Latency (ms)" },
      marks: [
        Plot.ruleY([xMax  ]),
        Plot.lineY(data, { x: d => new Date(d.date), y: "ms", stroke: "blue" })
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

export function BarGraph({ data }: any) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const dataMax = Math.max(...data.map((d: any) => d.column_value));
    const defaultMax = 500;

    const xMax = Math.max(defaultMax, dataMax);

    const chart = Plot.plot({
      width: chartRef.current.clientWidth,
      height: chartRef.current.clientHeight,
      marks: [
        Plot.axisX({ label: null, lineWidth: 8 }),
        Plot.axisY({ label: "Requests" }),
        Plot.ruleY([xMax]),
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

export function StackedBarsGraph({ data }: any) {
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

export function MultiLineChart({ data }: any) {
  const chartRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!chartRef.current) return;

    chartRef.current.innerHTML = "";

    const dataMax = Math.max(...data.map((d: any) => d.value));
    const defaultMax = 500;

    const xMax = Math.max(defaultMax, dataMax);

    const chart = Plot.plot({
      style: "overflow: visible;",
      y: { grid: true },
      width: chartRef.current.clientWidth,
      height: chartRef.current.clientHeight,
      marks: [
        Plot.ruleY([xMax]),
        Plot.lineY(data, { x: d => new Date(d.date), y: "value", stroke: "symbol" }),
        Plot.text(data, Plot.selectLast({ x: d => new Date(d.date), y: "value", z: "symbol", text: "symbol", textAnchor: "start", dx: 3 }))
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




