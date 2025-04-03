'use client';

import { useState } from 'react';
import axios from 'axios';
import { useRouter } from 'next/navigation';
import Navigation from "@/components/Navigation";

export default function CreateReportPage() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    timefrom: '',
    timeto: '',
    grouping: 'hour'
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      // Валидация дат
      if (!formData.timefrom || !formData.timeto) {
        throw new Error('Необходимо указать обе даты');
      }

      if (new Date(formData.timeto) < new Date(formData.timefrom)) {
        throw new Error('Дата окончания не может быть раньше даты начала');
      }

      const response = await axios.post(
        'http://localhost:3377/api/reports/generate',
        {
          timefrom: formData.timefrom,
          timeto: formData.timeto,
          grouping: formData.grouping
        },
        {
          withCredentials: true,
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );

      if (response.data.status === 'Error') {
        throw new Error(response.data.error?.detail || 'Ошибка генерации отчета');
      }

      // Перенаправляем на страницу отчетов после успешного создания
      router.push('/reports');
    } catch (err) {
      setError(
        axios.isAxiosError(err)
          ? err.response?.data?.error?.detail || 'Ошибка при создании отчета'
          : err instanceof Error
            ? err.message
            : 'Неизвестная ошибка'
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-md mx-auto bg-white p-8 rounded-lg shadow">
          <h1 className="text-2xl font-bold mb-6">Создать новый отчет</h1>
          
          {error && (
            <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="timefrom" className="block text-sm font-medium text-gray-700 mb-1">
                Дата начала
              </label>
              <input
                type="datetime-local"
                id="timefrom"
                name="timefrom"
                required
                value={formData.timefrom}
                onChange={handleChange}
                className="w-full p-2 border rounded"
              />
            </div>

            <div>
              <label htmlFor="timeto" className="block text-sm font-medium text-gray-700 mb-1">
                Дата окончания
              </label>
              <input
                type="datetime-local"
                id="timeto"
                name="timeto"
                required
                value={formData.timeto}
                onChange={handleChange}
                className="w-full p-2 border rounded"
              />
            </div>

            <div>
              <label htmlFor="grouping" className="block text-sm font-medium text-gray-700 mb-1">
                Группировка данных
              </label>
              <select
                id="grouping"
                name="grouping"
                value={formData.grouping}
                onChange={handleChange}
                className="w-full p-2 border rounded"
              >
                <option value="day">По дням</option>
                <option value="hour">По часам</option>
                <option value="minute">По минутам</option>
                <option value="second">По секундам</option>
              </select>
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={() => router.push('/reports')}
                className="px-4 py-2 border border-gray-300 rounded text-gray-700 hover:bg-gray-50"
              >
                Отмена
              </button>
              <button
                type="submit"
                disabled={loading}
                className={`px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 ${
                  loading ? 'opacity-50 cursor-not-allowed' : ''
                }`}
              >
                {loading ? 'Генерация...' : 'Генерировать отчет'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}