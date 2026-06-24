import type { ImageProps } from 'next/image'

export default function MockImage({ src, alt, width, height, ...props }: ImageProps) {
  return (
    // eslint-disable-next-line @next/next/no-img-element
    <img
      src={typeof src === 'string' ? src : ''}
      alt={alt}
      width={width as number}
      height={height as number}
      {...props}
    />
  )
}
