'use client';

import Navigation from "@/components/Navigation";
import axios from "axios";
import { useEffect, useState } from "react";
import Link from "next/link";

interface Role {
    title: string;
}

interface User {
    id: string;
    email: string;
    creation_date: string;
    role: Role;
}

export default function UsersListPage() {
    const [usersList, setUsersList] = useState<Array<User>>([]);
    const [filteredUsers, setFilteredUsers] = useState<Array<User>>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    
    // Фильтры
    const [emailFilter, setEmailFilter] = useState("");
    const [roleFilter, setRoleFilter] = useState("");
    const [dateFromFilter, setDateFromFilter] = useState("");
    const [dateToFilter, setDateToFilter] = useState("");

    useEffect(() => {
        const fetchUsersList = async () => {
            try {
                const response = await axios.get("http://localhost:3377/api/user/list", {
                    withCredentials: true,
                });
                
                if (response.data.status == "Error") {
                    throw new Error(response.data.error?.detail || "Failed to fetch users")
                }
                
                setUsersList(response.data?.data || []);
                setFilteredUsers(response.data?.data || []);
            } catch (err) {
                setError(err instanceof Error ? err.message : "Failed to fetch user data")
            } finally {
                setLoading(false);
            }
        };
        fetchUsersList();
    }, []);

    // Применение фильтров
    useEffect(() => {
        let result = [...usersList];
        
        if (emailFilter) {
            result = result.filter(user => 
                user.email.toLowerCase().includes(emailFilter.toLowerCase())
            );
        }
        
        if (roleFilter) {
            result = result.filter(user => 
                user.role.title.toLowerCase().includes(roleFilter.toLowerCase())
            );
        }
        
        if (dateFromFilter) {
            const fromDate = new Date(dateFromFilter);
            result = result.filter(user => 
                new Date(user.creation_date) >= fromDate
            );
        }
        
        if (dateToFilter) {
            const toDate = new Date(dateToFilter);
            result = result.filter(user => 
                new Date(user.creation_date) <= toDate
            );
        }
        
        setFilteredUsers(result);
    }, [usersList, emailFilter, roleFilter, dateFromFilter, dateToFilter]);

    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString("ru-RU", {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    };

    if (loading) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="min-h-screen flex items-center justify-center">
                <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
                    {error}
                </div>
            </div>
        );
    }

    if (!usersList.length) {
        return (
            <div className="min-h-screen">
                <Navigation />
                <div className="p-8">
                    <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded">
                        No users found
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen">
            <Navigation />
            
            <div className="p-8">
                <div className="flex justify-between items-center mb-6">
                    <h1 className="text-2xl font-bold">Список пользователей</h1>
                    <Link 
                        href="/users/create"
                        className="px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700"
                    >
                        Создать пользователя
                    </Link>
                </div>

                {/* Фильтры */}
                <div className="bg-white p-4 rounded-lg shadow mb-6">
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Поиск по email</label>
                            <input
                                type="text"
                                value={emailFilter}
                                onChange={(e) => setEmailFilter(e.target.value)}
                                className="w-full p-2 border rounded"
                                placeholder="Введите email"
                            />
                        </div>
                        
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Роль</label>
                            <select
                                value={roleFilter}
                                onChange={(e) => setRoleFilter(e.target.value)}
                                className="w-full p-2 border rounded"
                            >
                                <option value="">Все роли</option>
                                <option value="admin">Admin</option>
                                <option value="employee">Employee</option>
                            </select>
                        </div>
                        
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Дата с</label>
                            <input
                                type="date"
                                value={dateFromFilter}
                                onChange={(e) => setDateFromFilter(e.target.value)}
                                className="w-full p-2 border rounded"
                            />
                        </div>
                        
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">Дата по</label>
                            <input
                                type="date"
                                value={dateToFilter}
                                onChange={(e) => setDateToFilter(e.target.value)}
                                className="w-full p-2 border rounded"
                            />
                        </div>
                    </div>
                </div>

                {/* Таблица пользователей */}
                <div className="bg-white shadow rounded-lg overflow-hidden">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Роль</th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Дата создания</th>
                                {/* <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Действия</th> */}
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {filteredUsers.map((user) => (
                                <tr key={user.id}>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{user.email}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 capitalize">{user.role.title.toLowerCase()}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{formatDate(user.creation_date)}</td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                                        {/* <Link
                                            href={`/users/edit/${user.id}`}
                                            className="text-indigo-600 hover:text-indigo-900 mr-4"
                                        >
                                            Редактировать
                                        </Link> */}
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}