'use client';

import Link from "next/link";
import { useRouter } from "next/navigation";
import Navigation from "@/components/Navigation";

export default function MonitoringSelectionPage() {
    const router = useRouter();

    // Параметры мониторинга
    const monitoringOptions = [
        { id: 'temperature', name: 'Температура', color: 'bg-red-100 hover:bg-red-200', textColor: 'text-red-800' },
        { id: 'pressure', name: 'Давление', color: 'bg-blue-100 hover:bg-blue-200', textColor: 'text-blue-800' },
        { id: 'flap', name: 'Заслонки', color: 'bg-green-100 hover:bg-green-200', textColor: 'text-green-800' }
    ];

    return (
        <div className="min-h-screen bg-gray-50">
            <Navigation />

            <div className="container mx-auto px-4 py-12">
                <div>
                    <div className="text-center mb-12">
                        <h1 className="text-3xl font-bold text-gray-900 mb-2">Мониторинг в реальном времени</h1>
                        <p className="text-lg text-gray-600">Выберите область для отслеживания в реальном времени</p>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
                        {monitoringOptions.map((option) => (
                            <Link
                                key={option.id}
                                href={`/monitoring/realtime/${option.id}`}
                                className={`${option.color} ${option.textColor} rounded-xl shadow-md p-8 transition-all duration-300 transform hover:scale-105 hover:shadow-lg flex flex-col items-center`}
                            >
                                <div className="text-4xl mb-4">
                                    {option.id === 'temperature' && (
                                        <span className="text-red-500">📊</span>
                                    )}
                                    {option.id === 'pressure' && (
                                        <span className="text-blue-500">📊</span>
                                    )}
                                    {option.id === 'flap' && (
                                        <span className="text-green-500">📊</span>
                                    )}
                                </div>
                                <h2 className="text-xl font-semibold mb-2">{option.name}</h2>
                                <div className="mt-4 px-4 py-2 bg-white rounded-full text-xs font-medium shadow-inner">
                                    Отслеживать
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
                <div>
                    <div className="text-center mb-12 m-10">
                        <h1 className="text-3xl font-bold text-gray-900 mb-2">Поиск по данным мониторинга</h1>
                        <p className="text-lg text-gray-600">Выберите область для поиска</p>
                    </div>
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-4xl mx-auto">
                        {monitoringOptions.map((option) => (
                            <Link
                                key={option.id}
                                href={`/monitoring/search/${option.id}`}
                                className={`${option.color} ${option.textColor} rounded-xl shadow-md p-8 transition-all duration-300 transform hover:scale-105 hover:shadow-lg flex flex-col items-center`}
                            >
                                <div className="text-4xl mb-4">
                                    {option.id === 'temperature' && (
                                        <span className="text-red-500">📊</span>
                                    )}
                                    {option.id === 'pressure' && (
                                        <span className="text-blue-500">📊</span>
                                    )}
                                    {option.id === 'flap' && (
                                        <span className="text-green-500">📊</span>
                                    )}
                                </div>
                                <h2 className="text-xl font-semibold mb-2">{option.name}</h2>
                                <div className="mt-4 px-4 py-2 bg-white rounded-full text-xs font-medium shadow-inner">
                                    Искать
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}