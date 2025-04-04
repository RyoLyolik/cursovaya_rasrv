'use client';

import { useEffect, useState, useRef } from 'react';
import { Chart, registerables } from 'chart.js';
import Navigation from "@/components/Navigation";

Chart.register(...registerables);

// Желаемые значения температуры для каждой позиции (1-46)
const DESIRE_VALUES = [
  50, 60, 75, 100,
  125, 160, 200, 250,
  310, 400, 490, 570,
  630, 690, 740, 770,
  795, 815, 840, 855,
  870, 890, 910, 925,
  930, 940, 937, 930,
  915, 910, 890, 875,
  860, 845, 830, 800,
  760, 730, 650, 425,
  300, 225, 150, 120,
  80, 60,
];
const DEVIATION = 5

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

interface Alert {
  id: number;
  pos_id: number;
  currentValue: number;
  desiredValue: number;
  deviation: number;
  percentDeviation: number;
  timestamp: string;
}

export default function MonitoringPage() {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [readData, setReadData] = useState<PositionData>({});
  const [alerts, setAlerts] = useState<Alert[]>([]);
  const nextAlertId = useRef(1);
  
  const dataChartRef = useRef<Chart | null>(null);
  const dataCanvasRef = useRef<HTMLCanvasElement>(null);

  // Инициализация пустых данных для всех позиций
  const initializeData = () => {
    const emptyData: PositionData = {};
    for (let i = 1; i <= 46; i++) {
      emptyData[i] = null;
    }
    setReadData({...emptyData});
  };

  // Проверка отклонений температуры
  const checkDeviation = (pos_id: number, currentValue: number, ts: string) => {
    const desiredValue = DESIRE_VALUES[pos_id - 1];
    const deviation = currentValue - desiredValue;
    const percentDeviation = (Math.abs(deviation) / desiredValue) * 100;
    // Если отклонение больше 15%
    if (percentDeviation > DEVIATION) {
      const newAlert = {
        id: nextAlertId.current++,
        pos_id,
        currentValue,
        desiredValue,
        deviation,
        percentDeviation,
        timestamp: ts,
      };
      
      setAlerts(prev => [...prev, newAlert]);
    }
  };

  // Закрытие уведомления
  const closeAlert = (id: number) => {
    setAlerts(prev => prev.filter(alert => alert.id !== id));
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
            beginAtZero: false,
            title: {
              display: true,
              text: 'Значения'
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
    
    const values = Array.from({length: 46}, (_, i) => 
      data[i + 1]?.value || 0
    );
    
    chartRef.current.data.datasets[0].data = values;
    chartRef.current.update();
  };

  useEffect(() => {
    initializeData();

    // Инициализация графиков
    if (dataCanvasRef.current) {
      initChart(dataCanvasRef, dataChartRef, 'Значения температуры', 'rgba(255, 99, 132, 0.2)', 46);
    }

    // Подключение к WebSocket
    const ws = new WebSocket('ws://localhost:2114/ws');
    setSocket(ws);

    ws.onmessage = (event) => {
      const data: MonitoringData = JSON.parse(event.data);
      
      // Обновляем данные в соответствии с типом
      if (data.record_type === 'temperature') {
        setReadData(prev => {
          const newData = {...prev, [data.pos_id]: data};
          updateChart(dataChartRef, newData);
          
          // Проверяем отклонение температуры
          checkDeviation(data.pos_id, data.value, data.timestamp);
          return newData;
        });
      }
    };

    return () => {
      ws.close();
      if (dataChartRef.current) dataChartRef.current.destroy();
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
              <th className="py-2 px-4 border">Желаемое</th>
              <th className="py-2 px-4 border">Отклонение</th>
              <th className="py-2 px-4 border">Время</th>
            </tr>
          </thead>
          <tbody>
            {Array.from({length: leng}, (_, i) => i + 1).map(pos => (
              <tr key={pos} className={`hover:bg-gray-50 ${
                readData[pos]?.value && 
                Math.abs((readData[pos]?.value! - DESIRE_VALUES[pos-1]) / DESIRE_VALUES[pos-1]) > DEVIATION/100
                  ? 'bg-red-50' 
                  : ''
              }`}>
                <td className="py-2 px-4 border text-center">{pos}</td>
                <td className="py-2 px-4 border text-center">
                  {data[pos]?.value.toFixed(1) || '-'}
                </td>
                <td className="py-2 px-4 border text-center">
                  {DESIRE_VALUES[pos-1]}
                </td>
                <td className="py-2 px-4 border text-center">
                  {data[pos]?.value 
                    ? `${((data[pos]?.value! - DESIRE_VALUES[pos-1]) / DESIRE_VALUES[pos-1] * 100).toFixed(1)}%` 
                    : '-'}
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
        {/* Уведомления об отклонениях */}
        <div className="fixed top-4 right-4 space-y-2 z-50">
          {alerts.map(alert => (
            <div key={alert.id} className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded shadow-lg max-w-xs">
              <div className="flex justify-between items-start">
                <div>
                  <p className="font-bold">Отклонение от кривой обжига!</p>
                  <p className="text-sm">Позиция: {alert.pos_id}</p>
                  <p className="text-sm">Текущая: {alert.currentValue.toFixed(1)}</p>
                  <p className="text-sm">Желаемая: {alert.desiredValue}</p>
                  <p className="text-sm">Отклонение: {alert.deviation.toFixed(1)} ({alert.percentDeviation.toFixed(1)}%)</p>
                  <p className="text-sm">Время: {alert.timestamp}</p>
                </div>
                <button 
                  onClick={() => closeAlert(alert.id)}
                  className="text-red-500 hover:text-red-700"
                >
                  ×
                </button>
              </div>
            </div>
          ))}
        </div>

        <h1 className="text-2xl font-bold mb-6">Параметры температуры в печи</h1>
        
        {/* Графики */}
        <div className="gap-6 mb-8">
          <div className="bg-white p-4 rounded-lg shadow">
            <canvas ref={dataCanvasRef} height="100"></canvas>
          </div>
        </div>

        {/* Таблицы данных */}
        <div className="space-y-8">
          <DataTable title="Таблица параметров" data={readData} leng={46} />
        </div>
      </div>
    </div>
  );
}