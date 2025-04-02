'use client';

import { useState } from 'react';
import axios from 'axios';
import Navigation from "@/components/Navigation";

interface TemperatureRecord {
  position: number;
  value: number;
  timestamp: string;
}

export default function TemperatureHistoryPage() {
  const [records, setRecords] = useState<TemperatureRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Фильтры
  const [filters, setFilters] = useState({
    position: '',
    valueFrom: '',
    valueTo: '',
    timeFrom: '',
    timeTo: ''
  });

  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFilters(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const fetchTemperatureData = async () => {
    setLoading(true);
    setError(null);
    
    try {
      // Подготовка параметров запроса
      const params = new URLSearchParams();
      
      if (filters.position) params.append('position', filters.position);
      if (filters.valueFrom) params.append('ge', filters.valueFrom);
      if (filters.valueTo) params.append('le', filters.valueTo);
      if (filters.timeFrom) params.append('timefrom', filters.timeFrom);
      if (filters.timeTo) params.append('timeto', filters.timeTo);

      const response = await axios.get(
        `http://localhost:3377/api/history?datatype=temperature&${params.toString()}`,
        { withCredentials: true }
      );

      setRecords(response.data);
    } catch (err) {
      setError(
        axios.isAxiosError(err) 
          ? err.response?.data?.message || 'Ошибка при загрузке данных'
          : 'Неизвестная ошибка'
      );
      setRecords([]);
    } finally {
      setLoading(false);
    }
  };

  const formatDateTime = (isoString: string) => {
    return new Date(isoString).toLocaleString('ru-RU');
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="container mx-auto px-4 py-8">
        <h1 className="text-2xl font-bold mb-6">История температурных показаний</h1>
        
        {/* Фильтры */}
        <div className="bg-white p-6 rounded-lg shadow-md mb-6">
          <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Позиция (1-46)</label>
              <select
                name="position"
                value={filters.position}
                onChange={handleFilterChange}
                className="w-full p-2 border rounded"
              >
                <option value="">Все</option>
                {Array.from({length: 46}, (_, i) => i + 1).map(pos => (
                  <option key={pos} value={pos}>{pos}</option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Значение от</label>
              <input
                type="number"
                name="valueFrom"
                value={filters.valueFrom}
                onChange={handleFilterChange}
                className="w-full p-2 border rounded"
                placeholder="Минимальное значение"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Значение до</label>
              <input
                type="number"
                name="valueTo"
                value={filters.valueTo}
                onChange={handleFilterChange}
                className="w-full p-2 border rounded"
                placeholder="Максимальное значение"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Время с</label>
              <input
                type="datetime-local"
                name="timeFrom"
                value={filters.timeFrom}
                onChange={handleFilterChange}
                className="w-full p-2 border rounded"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Время до</label>
              <input
                type="datetime-local"
                name="timeTo"
                value={filters.timeTo}
                onChange={handleFilterChange}
                className="w-full p-2 border rounded"
              />
            </div>
          </div>
          
          <div className="mt-4 flex justify-end">
            <button
              onClick={fetchTemperatureData}
              disabled={loading}
              className={`px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 ${
                loading ? 'opacity-50 cursor-not-allowed' : ''
              }`}
            >
              {loading ? 'Загрузка...' : 'Применить фильтры'}
            </button>
          </div>
        </div>

        {/* Таблица с данными */}
        <div className="bg-white rounded-lg shadow overflow-hidden">
          {loading ? (
            <div className="p-8 flex justify-center">
              <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-indigo-500"></div>
            </div>
          ) : error ? (
            <div className="p-4 bg-red-100 border-l-4 border-red-500 text-red-700">
              <p>{error}</p>
            </div>
          ) : records.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              {Object.values(filters).some(f => f) 
                ? 'Нет данных, соответствующих выбранным фильтрам' 
                : 'Примените фильтры для отображения данных'}
            </div>
          ) : (
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Позиция</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Значение</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Время измерения</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {records.map((record, index) => (
                  <tr key={index} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{record.position}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{record.value.toFixed(4)}</td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDateTime(record.timestamp)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
}