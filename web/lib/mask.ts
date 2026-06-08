export type MaskPoint = {
  x: number;
  y: number;
};

export function serializeMask(points: MaskPoint[]): string {
  return JSON.stringify(points);
}
