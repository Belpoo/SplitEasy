/**
 * Archivo: server.js
 * Propósito: Punto de entrada del microservicio Auth Service.
 * Decisiones: Se escoge Fastify por su desempeño asincrónico superior.
 * Se configuran plugins clave como fastify-cookie para la manipulación HttpOnly,
 * el requisito central de seguridad impuesto en la sesión.
 */

require('dotenv').config();
const fastify = require('fastify')({ logger: true });
const authController = require('./controllers/auth.controller');
const { initDb } = require('./db');

// Registro de plugins
fastify.register(require('@fastify/cors'),
 {
  origin: true, // Configurable en prod
  credentials: true
});

fastify.register(require('@fastify/cookie'), {
  secret: process.env.COOKIE_SECRET || "secreto-seguro-cookie",
  parseOptions: {} 
});

// Rutas de autenticación
fastify.post('/auth/register', authController.register);
fastify.post('/auth/login', authController.login);
fastify.post('/auth/refresh', authController.refresh);
fastify.post('/auth/logout', authController.logout);
fastify.get('/auth/me', authController.me);

const start = async () => {
  try {
    // Inicializar las tablas relacionales primero (si no existen)
    await initDb();
    
    // El port suele asignarse por Render mediante PORT
    const port = process.env.PORT || 3001;
    await fastify.listen({ port, host: '0.0.0.0' });
    fastify.log.info(`Auth Service corriendo en el puerto ${port}`);
  } catch (err) {
    fastify.log.error(err);
    process.exit(1);
  }
};

// Exportar fastify app para que en Jest podamos inyectar peticiones sin inicializar red
if (require.main === module) {
  start();
} else {
  module.exports = fastify;
}
