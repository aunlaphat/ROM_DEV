/**
 * Enhanced logger utility for React applications
 * Provides consistent formatting and categorized logging functionality
 */

// Types
type LogLevel = 'debug' | 'info' | 'warn' | 'error' | 'success';
type LogIcon = 'api' | 'auth' | 'login' | 'logout' | 'state' | 'time' | 'error' | 'success' | 'warning' | 'navigation';

// ในการกำหนด interface ที่มี key เป็น union type ใช้ Record แทน mapped type
type LogColors = Record<LogLevel, string>;
type LogIcons = Record<LogIcon, string>;

// Constants
const LOG_COLORS: LogColors = {
  debug: '#808080',   // Gray
  info: '#1890ff',    // Primary Blue
  warn: '#faad14',    // Warning Yellow
  error: '#f5222d',   // Error Red
  success: '#52c41a'  // Success Green
};

const LOG_ICONS: LogIcons = {
  api: '🌐',
  auth: '🔒',
  login: '🔑',
  logout: '👋',
  state: '📊',
  time: '⏱️',
  error: '❌',
  success: '✅',
  warning: '⚠️',
  navigation: '🔄'
};

/**
 * Gets formatted timestamp for logs
 * @returns Formatted time string
 */
const getTimestamp = (): string => {
  return new Date().toLocaleTimeString();
};

/**
 * Determines if we're in production environment
 * (Logging can be reduced or disabled in production)
 */
const isProduction = (): boolean => {
  return process.env.NODE_ENV === 'production';
};

/**
 * Main logger object
 */
export const logger = {
  /**
   * Base logging method
   * @param level - Severity level of the log
   * @param message - Message to log
   * @param data - Optional data to include
   */
  log: (level: LogLevel, message: string, data?: any): void => {
    // Skip debug logs in production
    if (isProduction() && level === 'debug') return;
    
    const color = LOG_COLORS[level];
    const timestamp = getTimestamp();
    
    console.log(
      `%c${timestamp} ${message}`, 
      `color: ${color}; font-weight: bold`,
      data !== undefined ? data : ''
    );
  },

  /**
   * Debug level logging - for development details
   */
  debug: (message: string, data?: any): void => {
    logger.log('debug', message, data);
  },

  /**
   * Info level logging - general information
   */
  info: (message: string, data?: any): void => {
    logger.log('info', message, data);
  },

  /**
   * Warning level logging
   */
  warn: (message: string, data?: any): void => {
    logger.log('warn', `${LOG_ICONS.warning} ${message}`, data);
  },

  /**
   * Error level logging
   */
  error: (message: string, error?: any): void => {
    logger.log('error', `${LOG_ICONS.error} ${message}`, error);
    
    // In development, we can also log to the console.error for better DevTools visibility
    if (!isProduction() && error) {
      console.error(error);
    }
  },

  /**
   * Success level logging
   */
  success: (message: string, data?: any): void => {
    logger.log('success', `${LOG_ICONS.success} ${message}`, data);
  },

  /**
   * API related logging
   */
  api: {
    request: (endpoint: string, data?: any): void => {
      logger.log('info', `${LOG_ICONS.api} API Request: ${endpoint}`, data);
    },
    
    success: (endpoint: string, data?: any): void => {
      logger.log('success', `${LOG_ICONS.api} API Success: ${endpoint}`, data);
    },
    
    error: (endpoint: string, error: any): void => {
      logger.log('error', `${LOG_ICONS.api} API Error: ${endpoint}`, error);
    }
  },

  /**
   * Authentication related logging
   */
  auth: {
    login: (username: string, data?: any): void => {
      logger.log('info', `${LOG_ICONS.login} Login: ${username}`, data);
    },
    
    logout: (username?: string, data?: any): void => {
      const userInfo = username ? `: ${username}` : '';
      logger.log('info', `${LOG_ICONS.logout} Logout${userInfo}`, data);
    },
    
    authCheck: (status: string, data?: any): void => {
      logger.log('info', `${LOG_ICONS.auth} Auth Check: ${status}`, data);
    },
    
    error: (message: string, error?: any): void => {
      logger.log('error', `${LOG_ICONS.auth} Auth Error: ${message}`, error);
    }
  },

  /**
   * Navigation related logging
   */
  navigation: {
    navigate: (destination: string, data?: any): void => {
      logger.log('info', `${LOG_ICONS.navigation} Navigate to: ${destination}`, data);
    },
    
    redirect: (from: string, to: string, data?: any): void => {
      logger.log('info', `${LOG_ICONS.navigation} Redirect: ${from} → ${to}`, data);
    }
  },

  /**
   * Performance monitoring utilities
   */
  perf: {
    start: (label: string): void => {
      if (isProduction()) return;
      
      console.group(`${LOG_ICONS.time} ${label}`);
      console.time(label);
    },
    
    end: (label: string): void => {
      if (isProduction()) return;
      
      console.timeEnd(label);
      console.groupEnd();
    },
    
    measure: async <T>(label: string, fn: () => Promise<T>): Promise<T> => {
      logger.perf.start(label);
      try {
        const result = await fn();
        return result;
      } finally {
        logger.perf.end(label);
      }
    }
  },

  /**
   * State management related logging
   */
  state: {
    update: (action: string, data?: any): void => {
      logger.log('info', `${LOG_ICONS.state} State Update: ${action}`, data);
    },
    
    error: (message: string, error?: any): void => {
      logger.log('error', `${LOG_ICONS.state} State Error: ${message}`, error);
    }
  }
};

// Export default for convenience
export default logger;