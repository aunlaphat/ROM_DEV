type LogLevel = 'debug' | 'info' | 'warn' | 'error';
type LogModule = 'auth' | 'api' | 'route' | 'state' | 'user' | 'error' | 'success' | 'db' | 'cache';

const LOG_COLORS = {
  debug: '#808080',  // gray
  info: '#0077FF',   // blue
  warn: '#FFA500',   // orange
  error: '#FF0000',  // red
  success: '#00FF00' // green
};

const LOG_EMOJI = {
  auth: '🔐',
  api: '🌐',
  route: '🔄',
  state: '📊',
  user: '👤',
  error: '❌',
  success: '✅',
  db: '💾',
  cache: '⚡',
  warning: '⚠️',
  time: '⏱️'
};

const formatTime = () => {
  return new Date().toISOString().split('T')[1].split('.')[0];
};

const createLogMessage = (
  module: LogModule, 
  level: LogLevel, 
  message: string, 
  data?: any
) => {
  const emoji = LOG_EMOJI[module];
  const color = LOG_COLORS[level];
  const time = formatTime();
  
  return {
    text: `%c${emoji} [${time}] [${module.toUpperCase()}] ${message}`,
    style: `color: ${color}; font-weight: bold`,
    data: data || ''
  };
};

export const logger = {
  auth: (level: LogLevel, message: string, data?: any) => {
    const log = createLogMessage('auth', level, message, data);
    console.log(log.text, log.style, log.data);
  },

  api: (level: LogLevel, message: string, data?: any) => {
    const log = createLogMessage('api', level, message, data);
    console.log(log.text, log.style, log.data);
  },

  route: (level: LogLevel, message: string, data?: any) => {
    const log = createLogMessage('route', level, message, data);
    console.log(log.text, log.style, log.data);
  },

  state: (level: LogLevel, message: string, data?: any) => {
    const log = createLogMessage('state', level, message, data);
    console.log(log.text, log.style, log.data);
  },

  user: (level: LogLevel, message: string, data?: any) => {
    const log = createLogMessage('user', level, message, data);
    console.log(log.text, log.style, log.data);
  },

  time: (label: string) => {
    console.time(`${LOG_EMOJI.time} ${label}`);
  },

  timeEnd: (label: string) => {
    console.timeEnd(`${LOG_EMOJI.time} ${label}`);
  },

  group: (label: string) => {
    console.group(`${LOG_EMOJI.success} ${label}`);
  },

  groupEnd: () => {
    console.groupEnd();
  },

  // Utility method for development only
  dev: (message: string, data?: any) => {
    if (process.env.NODE_ENV === 'development') {
      const log = createLogMessage('db', 'debug', `[DEV] ${message}`, data);
      console.log(log.text, log.style, log.data);
    }
  },

  // Method for critical errors that should always be logged
  critical: (message: string, error?: any) => {
    const log = createLogMessage('error', 'error', `[CRITICAL] ${message}`, error);
    console.error(log.text, log.style, log.data);
    // Could add error reporting service integration here
  }
};

// Usage examples:
// logger.api('info', 'API call started', { endpoint: '/users' });
// logger.time('API Call');
// logger.timeEnd('API Call');
// logger.group('User Authentication');
// logger.auth('info', 'User logged in');
// logger.groupEnd();
// logger.dev('Debug message', { data: 'test' });
// logger.critical('Database connection failed', error);
