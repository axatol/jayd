import { Link, LinkProps } from "react-router-dom";

export const ExternalLink = (props: LinkProps) => (
  <Link {...props} rel="noopener noreferrer" target="_blank" />
);

export const YoutubeChannelLink = (props: { id: string; name: string }) => (
  <ExternalLink
    rel="noopener noreferrer"
    target="_blank"
    to={`https://youtube.com/${props.id}`}
  >
    {props.name}
  </ExternalLink>
);

export const YoutubeVideoLink = (props: { id: string; title: string }) => (
  <ExternalLink
    rel="noopener noreferrer"
    target="_blank"
    to={`https://youtube.com/watch?v=${props.id}`}
  >
    {props.title}
  </ExternalLink>
);
