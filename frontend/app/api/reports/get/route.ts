import { NextRequest, NextResponse } from 'next/server';
import axios from 'axios';

export async function GET(request: NextRequest) {
  try {
    const { searchParams } = new URL(request.url);
    const filename = searchParams.get('filename');

    if (!filename) {
      return new NextResponse('Filename parameter is required', { status: 400 });
    }

    const response = await axios.get(
      `http://localhost:9999/api/report/file?filename=${encodeURIComponent(filename)}`,
      {
        responseType: 'arraybuffer', // Важно для бинарных данных
        headers: {
          'Accept': 'application/octet-stream',
        },
      }
    );
    // Создаем ответ с бинарными данными
    return new NextResponse(response.data, {
      status: 200,
      headers: new Headers({
        'Content-Disposition': response.headers['content-disposition'],
        'Content-Type': response.headers['content-type'] || 'application/octet-stream',
        'Content-Length': response.headers['content-length'] || String(response.data.length),
      }),
    });
  } catch (err) {
    if (axios.isAxiosError(err)) {
      return new NextResponse(JSON.stringify(err.response?.data || { status: 'Error' }), {
        status: err.response?.status || 500,
        headers: { 'Content-Type': 'application/json' },
      });
    }
    return new NextResponse(JSON.stringify({ status: 'Error' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' },
    });
  }
}