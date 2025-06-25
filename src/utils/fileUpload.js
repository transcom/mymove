export default function appendTimestampToFilename(file) {
  const now = new Date();
  const timestamp = now
    .toISOString()
    .replace(/[-T:.Z]/g, '')
    .slice(0, 14);
  const [name, extension = ''] = file.name?.split(/\.(?=[^.]+$)/) ?? '';
  const newFileName = `${name}-${timestamp}${extension ? `.${extension}` : ''}`;
  return new File([file], newFileName, { type: file.type });
}
