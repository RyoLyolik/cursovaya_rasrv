'use server';

import axios from "axios";
import { createSearchParamsFromClient } from "next/dist/server/request/search-params";
import { NextRequest, NextResponse } from "next/server";
import { cookies } from 'next/headers'


export async function POST(request: NextRequest) {
  const formData = await request.json();
  try {
    const resp = await axios.post('http://localhost:9999/api/report/generate', formData);
    
    // Создаем ответ
    const response = NextResponse.json({
      status: "Ok"
    });
    return response;
  } catch (err) {
    if (axios.isAxiosError(err)) {
      const response = NextResponse.json(err.response?.data || {status: 'Error'});
      return response
    } else {
      return NextResponse.json({status: 'Error'});
    }
  }
}