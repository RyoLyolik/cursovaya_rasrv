import { NextResponse } from 'next/server';
import axios from 'axios';

export async function GET(request: Request) {
  try {
    // Получаем параметры запроса из URL
    const { searchParams } = new URL(request.url);
    
    // Подготавливаем параметры для бэкенда
    const backendParams = new URLSearchParams();
    
    // Переносим поддерживаемые параметры
    const supportedParams = [
      'datatype',
      'position',
      'ge',
      'le', 
      'timefrom', 
      'timeto'
    ];
    
    supportedParams.forEach(param => {
      const value = searchParams.get(param);
      if (value) {
        backendParams.append(param, value);
      }
    });

    // Делаем запрос к бэкенду
    const backendUrl = `http://localhost:9999/api/data/find?${backendParams.toString()}`;
    
    const response = await axios.get(backendUrl, {
      // Если нужно передать куки или заголовки авторизации
      withCredentials: true,
      headers: {
        // Можно добавить необходимые заголовки
        'Accept': 'application/json',
      }
    });

    // Возвращаем данные от бэкенда
    return NextResponse.json(response.data);
    
  } catch (error) {
    // Обработка ошибок
    if (axios.isAxiosError(error)) {
      return NextResponse.json(
        { 
          status: 'Error',
          error: {
            message: error.response?.data?.message || 'Backend request failed',
            detail: error.response?.data?.error || error.message
          }
        },
        { status: error.response?.status || 500 }
      );
    }
    
    return NextResponse.json(
      { 
        status: 'Error',
        error: {
          message: 'Internal server error',
          detail: 'Unknown error occurred'
        }
      },
      { status: 500 }
    );
  }
}