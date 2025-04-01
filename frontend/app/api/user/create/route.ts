'use server';

import axios from "axios";
import { createSearchParamsFromClient } from "next/dist/server/request/search-params";
import { NextRequest, NextResponse } from "next/server";
import { cookies } from 'next/headers'


export async function POST(request: NextRequest) {
  const formData = await request.json();
  const sessionCookie = request.cookies.get('brick-session-id')
  try {
    const resp = await axios.post(
        'http://localhost:8080/account',
        formData,
        {
            headers: {
                Cookie: `${sessionCookie?.name}=${sessionCookie?.value}`
            },
            withCredentials: true
        }
    );
    
    // Создаем ответ
    const response = NextResponse.json({
      status: "Ok"
    });
    
    // Если есть куки в ответе от бэкенда, устанавливаем их
    return response;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      const response = NextResponse.json(err.response?.data || {status: 'Error'});
    } else {
      return NextResponse.json({status: 'Error'});
    }
  }
}