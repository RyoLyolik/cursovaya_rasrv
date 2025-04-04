'use client';

import { useEffect, useState, useRef } from 'react';
import { Chart, registerables } from 'chart.js';
import Navigation from "@/components/Navigation";

Chart.register(...registerables);

interface MonitoringData {
  timestamp: string;
  record_id: number;
  record_type: 'flap' | 'temperature' | 'pressure';
  value: number;
  pos_id: number;
}

interface PositionData {
  [key: number]: MonitoringData | null;
}

export default function MonitoringPage() {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [flapData, setflapData] = useState<PositionData>({});
  
  const flapChartRef = useRef<Chart | null>(null);
  
  const flapCanvasRef = useRef<HTMLCanvasElement>(null);

  // Инициализация пустых данных для всех позиций
  const initializeData = () => {
    const emptyData: PositionData = {};
    for (let i = 1; i <= 10; i++) {
      emptyData[i] = null;
    }
    setflapData({...emptyData});
  };

  // Инициализация графиков
  const initChart = (
    canvasRef: React.RefObject<HTMLCanvasElement>,
    chartRef: React.MutableRefObject<Chart | null>,
    label: string,
    color: string,
    leng: number
  ) => {
    if (!canvasRef.current) return;
    
    if (chartRef.current) {
      chartRef.current.destroy();
    }

    const ctx = canvasRef.current.getContext('2d');
    if (!ctx) return;

    chartRef.current = new Chart(ctx, {
      type: 'bar',
      data: {
        labels: Array.from({length: leng}, (_, i) => (i + 1).toString()),
        datasets: [{
          label: label,
          data: Array(leng).fill(0),
          backgroundColor: color,
          borderColor: color.replace('0.2', '1'),
          borderWidth: 1
        }]
      },
      options: {
        responsive: true,
        scales: {
          y: {
            min: 0,
            max: 100,
            title: {
              display: true,
              text: 'Значение'
            }
          },
          x: {
            title: {
              display: true,
              text: 'Позиция'
            }
          }
        },
        plugins: {
          title: {
            display: true,
            text: `${label} по позициям`
          }
        }
      }
    });
  };

  // Обновление данных на графике
  const updateChart = (
    chartRef: React.MutableRefObject<Chart | null>, 
    data: PositionData
  ) => {
    if (!chartRef.current) return;
    
    const values = Array.from({length: 10}, (_, i) => 
      data[i + 1]?.value || 0
    );
    
    chartRef.current.data.datasets[0].data = values;
    chartRef.current.update();
  };

  useEffect(() => {
    initializeData();

    // Инициализация графиков
    if (flapCanvasRef.current) {
      initChart(flapCanvasRef, flapChartRef, 'Значения положения заслонок', 'rgba(255, 99, 132, 0.2)', 10);
    }

    // Подключение к WebSocket
    const ws = new WebSocket('ws://localhost:2114/ws');
    setSocket(ws);

    ws.onmessage = (event) => {
      const data: MonitoringData = JSON.parse(event.data);
      
      // Обновляем данные в соответствии с типом
      switch(data.record_type) {
        case 'flap':
          setflapData(prev => {
            const newData = {...prev, [data.pos_id]: data};
            updateChart(flapChartRef, newData);
            return newData;
          });
          break;
      }
    };

    return () => {
      ws.close();
      if (flapChartRef.current) flapChartRef.current.destroy();
    };
  }, []);

  // Форматирование даты для отображения
  const formatDate = (dateString: string) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString();
  };

  // Компонент таблицы
  const DataTable = ({ title, data, leng }: { title: string; data: PositionData, leng: number }) => (
    <div className="mb-8">
      <h3 className="text-lg font-semibold mb-2">{title}</h3>
      <div className="overflow-x-auto">
        <table className="min-w-full bg-white border">
          <thead>
            <tr className="bg-gray-100">
              <th className="py-2 px-4 border">Позиция</th>
              <th className="py-2 px-4 border">Значение</th>
              <th className="py-2 px-4 border">Время</th>
            </tr>
          </thead>
          <tbody>
            {Array.from({length: leng}, (_, i) => i + 1).map(pos => (
              <tr key={pos} className="hover:bg-gray-50">
                <td className="py-2 px-4 border text-center">{pos}</td>
                <td className="py-2 px-4 border text-center">
                  {data[pos]?.value.toFixed(4) || '-'}
                </td>
                <td className="py-2 px-4 border text-center">
                  {formatDate(data[pos]?.timestamp)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold mb-6">Положение заслонок в печи</h1>
        
        {/* Графики */}
        <div className="gap-6 mb-8">
          <div className="bg-white p-4 rounded-lg shadow">
            <canvas ref={flapCanvasRef} height="100"></canvas>
          </div>
        </div>

        {/* Таблицы данных */}
        <div className="space-y-8">
          <DataTable title="Таблица параметров" data={flapData} leng={10} />
        </div>
      </div>
    </div>
  );
}