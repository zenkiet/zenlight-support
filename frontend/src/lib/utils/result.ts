export type Result<T, E = string> =
| { success: true; data: T }
| { success: false; error: E };

export async function safeCall<T>(promise: Promise<T>): Promise<Result<T>> {
  try {
    const data = await promise;
    return { success: true, data };
  } catch (e) {
    return { success: false, error: String(e) };
  }
}