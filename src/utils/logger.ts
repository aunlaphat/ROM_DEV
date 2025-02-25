type LogLevel = 'debug' | 'info' | 'warn' | 'error';

const LOG_COLORS = {
  debug: '#808080', // gray
  info: '#0077FF',  // blue
  warn: '#FFA500',  // orange
  error: '#FF0000', // red
};

const LOG_EMOJI = {
  auth: '🔐',
  api: '🌐',
  route: '🔄',
  state: '📊',
  user: '👤',
  error: '❌',
  success: '✅',
};

export const logger = {
  auth: (level: LogLevel, message: string, data?: any) => {
    const emoji = LOG_EMOJI.auth;
    const color = LOG_COLORS[level];
    console.log(
      `%c${emoji} [Auth] ${message}`,
      `color: ${color}; font-weight: bold`,
      data ? data : ''
    );
  },

  api: (level: LogLevel, message: string, data?: any) => {
    const emoji = LOG_EMOJI.api;
    const color = LOG_COLORS[level];
    console.log(
      `%c${emoji} [API] ${message}`,
      `color: ${color}; font-weight: bold`,
      data ? data : ''
    );
  },

  route: (level: LogLevel, message: string, data?: any) => {
    const emoji = LOG_EMOJI.route;
    const color = LOG_COLORS[level];
    console.log(
      `%c${emoji} [Route] ${message}`,
      `color: ${color}; font-weight: bold`,
      data ? data : ''
    );
  }
};
